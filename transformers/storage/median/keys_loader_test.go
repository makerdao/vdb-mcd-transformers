package median_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/median"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities/wards"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Median storage keys loader", func() {
	var (
		storageRepository *test_helpers.MockMakerStorageRepository
		storageKeysLoader storage.KeysLoader
	)

	BeforeEach(func() {
		storageRepository = &test_helpers.MockMakerStorageRepository{}
		storageKeysLoader = median.NewKeysLoader(storageRepository, test_data.EthMedianAddress())
	})

	It("returns value metadata for static keys", func() {
		mappings, err := storageKeysLoader.LoadMappings()
		Expect(err).NotTo(HaveOccurred())

		Expect(mappings[median.ValAndAgeStorageKey]).To(Equal(types.ValueMetadata{
			Name:        mcdStorage.Packed,
			Type:        types.PackedSlot,
			PackedTypes: map[int]types.ValueType{0: types.Uint128, 1: types.Uint32},
			PackedNames: map[int]string{0: median.Val, 1: median.Age},
		}))

		Expect(mappings[median.BarKey]).To(Equal(median.BarMetadata))
	})

	Describe("wards", func() {
		It("returns value metadata for wards", func() {
			wardsUser := fakes.FakeAddress.Hex()
			storageRepository.WardsKeys = []string{wardsUser}
			paddedWardsUser := "0x000000000000000000000000" + wardsUser[2:]
			wardsKey := common.BytesToHash(crypto.Keccak256(common.FromHex(paddedWardsUser + wards.WardsMappingIndex)))
			expectedMetadata := types.ValueMetadata{
				Name: wards.Wards,
				Keys: map[types.Key]string{constants.User: fakes.FakeAddress.Hex()},
				Type: types.Uint256,
			}

			mappings, err := storageKeysLoader.LoadMappings()
			Expect(err).NotTo(HaveOccurred())

			Expect(storageRepository.GetWardsKeysCalledWith).To(Equal(test_data.EthMedianAddress()))
			Expect(mappings[wardsKey]).To(Equal(expectedMetadata))
		})

		It("returns error on failure", func() {
			storageRepository.GetWardsKeysError = fakes.FakeError

			_, err := storageKeysLoader.LoadMappings()

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})
	})

	Describe("bud", func() {
		It("returns value metadata for bud", func() {
			budAddress := common.HexToAddress(test_data.RandomString(40)).Hex()
			storageRepository.MedianBudAddresses = []string{budAddress}
			paddedBudAddress := "0x000000000000000000000000" + budAddress[2:]
			budKey := common.BytesToHash(crypto.Keccak256(common.FromHex(paddedBudAddress + median.BudMappingIndex)))
			expectedMetadata := types.ValueMetadata{
				Name: median.Bud,
				Keys: map[types.Key]string{constants.A: budAddress},
				Type: types.Uint256,
			}

			mappings, err := storageKeysLoader.LoadMappings()
			Expect(err).NotTo(HaveOccurred())

			Expect(storageRepository.GetMedianBudAddressesCalledWith).To(Equal(test_data.EthMedianAddress()))
			Expect(mappings[budKey]).To(Equal(expectedMetadata))
		})

		It("returns error on failure", func() {
			storageRepository.GetMedianBudAddressesError = fakes.FakeError

			_, err := storageKeysLoader.LoadMappings()

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})
	})
})
