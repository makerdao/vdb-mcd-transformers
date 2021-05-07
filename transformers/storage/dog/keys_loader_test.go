package dog_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/dog"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities/wards"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	vdbStorage "github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Dog storage keys loader", func() {
	var (
		storageRepository *test_helpers.MockMakerStorageRepository
		storageKeysLoader storage.KeysLoader
	)

	BeforeEach(func() {
		storageRepository = &test_helpers.MockMakerStorageRepository{}
		storageKeysLoader = dog.NewKeysLoader(storageRepository, test_data.Dog130Address())
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

			Expect(storageRepository.GetWardsKeysCalledWith).To(Equal(test_data.Dog130Address()))
			Expect(mappings[wardsKey]).To(Equal(expectedMetadata))
		})

		Describe("when getting users fails", func() {
			It("returns error", func() {
				storageRepository.GetWardsKeysError = fakes.FakeError

				_, err := storageKeysLoader.LoadMappings()

				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(fakes.FakeError))
			})
		})
	})

	Describe("ilks", func() {
		Describe("when getting ilks fails", func() {
			It("returns error", func() {
				storageRepository.GetIlksError = fakes.FakeError

				_, err := storageKeysLoader.LoadMappings()

				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(fakes.FakeError))
			})
		})

		Describe("when getting ilks succeeds", func() {
			var (
				ilkClipKey = common.BytesToHash(crypto.Keccak256(common.FromHex(test_helpers.FakeIlk + dog.IlksMappingIndex)))
				mappings   map[common.Hash]types.ValueMetadata
			)

			BeforeEach(func() {
				storageRepository.Ilks = []string{test_helpers.FakeIlk}
				var err error
				mappings, err = storageKeysLoader.LoadMappings()
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns value metadata for ilk clip", func() {
				expectedMetadata := types.ValueMetadata{
					Name: dog.IlkClip,
					Keys: map[types.Key]string{constants.Ilk: test_helpers.FakeIlk},
					Type: types.Address,
				}

				Expect(mappings[ilkClipKey]).To(Equal(expectedMetadata))
			})

			It("returns value metadata for ilk chop", func() {
				ilkChopKey := vdbStorage.GetIncrementedKey(ilkClipKey, 1)
				expectedMetadata := types.ValueMetadata{
					Name: dog.IlkChop,
					Keys: map[types.Key]string{constants.Ilk: test_helpers.FakeIlk},
					Type: types.Uint256,
				}

				Expect(mappings[ilkChopKey]).To(Equal(expectedMetadata))
			})
		})
	})

	It("returns storage value metadata for static keys", func() {
		mappings, err := storageKeysLoader.LoadMappings()
		Expect(err).NotTo(HaveOccurred())

		Expect(mappings[dog.VatStorageKey]).To(Equal(types.ValueMetadata{
			Name: dog.Vat,
			Keys: nil,
			Type: types.Address,
		}))

		Expect(mappings[dog.VowStorageKey]).To(Equal(types.ValueMetadata{
			Name: dog.Vow,
			Keys: nil,
			Type: types.Address,
		}))

		Expect(mappings[dog.LiveStorageKey]).To(Equal(types.ValueMetadata{
			Name: dog.Live,
			Keys: nil,
			Type: types.Uint256,
		}))

		Expect(mappings[dog.HoleStorageKey]).To(Equal(types.ValueMetadata{
			Name: dog.Hole,
			Keys: nil,
			Type: types.Uint256,
		}))

		Expect(mappings[dog.DirtStorageKey]).To(Equal(types.ValueMetadata{
			Name: dog.Dirt,
			Keys: nil,
			Type: types.Uint256,
		}))
	})
})
