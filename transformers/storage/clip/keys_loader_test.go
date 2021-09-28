package clip_test

import (
	"math/rand"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities/wards"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	vdbStorage "github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Clip storage keys loader", func() {
	var (
		storageRepository *test_helpers.MockMakerStorageRepository
		storageKeysLoader storage.KeysLoader
	)

	BeforeEach(func() {
		storageRepository = &test_helpers.MockMakerStorageRepository{}
		storageKeysLoader = clip.NewKeysLoader(storageRepository, test_data.ClipLinkAV130Address())
	})

	It("returns value metadata for static keys", func() {
		mappings, err := storageKeysLoader.LoadMappings()

		Expect(err).NotTo(HaveOccurred())
		Expect(mappings[clip.DogKey]).To(Equal(clip.DogMetadata))
		Expect(mappings[clip.VowKey]).To(Equal(clip.VowMetadata))
		Expect(mappings[clip.SpotterKey]).To(Equal(clip.SpotterMetadata))
		Expect(mappings[clip.CalcKey]).To(Equal(clip.CalcMetadata))
		Expect(mappings[clip.BufKey]).To(Equal(clip.BufMetadata))
		Expect(mappings[clip.TailKey]).To(Equal(clip.TailMetadata))
		Expect(mappings[clip.CuspKey]).To(Equal(clip.CuspMetadata))
		Expect(mappings[clip.ChipAndTipStorageKey]).To(Equal(types.ValueMetadata{
			Name:        mcdStorage.Packed,
			Type:        types.PackedSlot,
			PackedTypes: map[int]types.ValueType{0: types.Uint64, 1: types.Uint192},
			PackedNames: map[int]string{0: clip.Chip, 1: clip.Tip},
		}))
		Expect(mappings[clip.ChostKey]).To(Equal(clip.ChostMetadata))
		Expect(mappings[clip.KicksKey]).To(Equal(clip.KicksMetadata))
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

			Expect(storageRepository.GetWardsKeysCalledWith).To(Equal(test_data.ClipLinkAV130Address()))
			Expect(mappings[wardsKey]).To(Equal(expectedMetadata))
		})
	})

	Describe("sales", func() {
		var (
			storageRepository *test_helpers.MockMakerStorageRepository
			storageKeysLoader storage.KeysLoader
			fakeUint256       = strconv.Itoa(rand.Intn(1000000))
		)
		BeforeEach(func() {
			storageRepository = &test_helpers.MockMakerStorageRepository{}
			storageKeysLoader = clip.NewKeysLoader(storageRepository, test_data.ClipLinkAV130Address())
		})
		Describe("when getting sales fails", func() {
			It("returns an error", func() {
				storageRepository.GetClipSalesError = fakes.FakeError

				_, err := storageKeysLoader.LoadMappings()

				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(fakes.FakeError))
			})
		})

		Describe("when getting sales succeeds", func() {
			var (
				mappings          map[common.Hash]types.ValueMetadata
				fakeUint256Hex, _ = shared.ConvertIntStringToHex(fakeUint256)
				salesPosKey       = common.BytesToHash(crypto.Keccak256(common.FromHex(fakeUint256Hex + clip.SalesMappingIndex)))
			)

			BeforeEach(func() {
				storageRepository.ClipSalesIDs = []string{fakeUint256}
				var err error
				mappings, err = storageKeysLoader.LoadMappings()
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns value metadata for sales pos", func() {
				expectedMetadata := types.ValueMetadata{
					Name: clip.SalePos,
					Keys: map[types.Key]string{constants.SaleId: fakeUint256},
					Type: types.Uint256,
				}

				Expect(mappings[salesPosKey]).To(Equal(expectedMetadata))
			})

			It("returns value metadata for sales tab", func() {
				salesTabKey := vdbStorage.GetIncrementedKey(salesPosKey, 1)
				expectedMetadata := types.ValueMetadata{
					Name: clip.SaleTab,
					Keys: map[types.Key]string{constants.SaleId: fakeUint256},
					Type: types.Uint256,
				}

				Expect(mappings[salesTabKey]).To(Equal(expectedMetadata))
			})

			It("returns value metadata for sales lot", func() {
				salesLotKey := vdbStorage.GetIncrementedKey(salesPosKey, 2)
				expectedMetadata := types.ValueMetadata{
					Name: clip.SaleLot,
					Keys: map[types.Key]string{constants.SaleId: fakeUint256},
					Type: types.Uint256,
				}

				Expect(mappings[salesLotKey]).To(Equal(expectedMetadata))
			})

			It("returns value metadata for sale usr + tic packed slot", func() {
				saleUsrKey := vdbStorage.GetIncrementedKey(salesPosKey, 3)
				expectedMetadata := types.ValueMetadata{
					Name:        mcdStorage.Packed,
					Keys:        map[types.Key]string{constants.SaleId: fakeUint256},
					Type:        types.PackedSlot,
					PackedTypes: map[int]types.ValueType{0: types.Address, 1: types.Uint96},
					PackedNames: map[int]string{0: clip.SaleUsr, 1: clip.SaleTic},
				}

				Expect(mappings[saleUsrKey]).To(Equal(expectedMetadata))
			})

			It("returns value metadata for sales top", func() {
				salesTopKey := vdbStorage.GetIncrementedKey(salesPosKey, 4)
				expectedMetadata := types.ValueMetadata{
					Name: clip.SaleTop,
					Keys: map[types.Key]string{constants.SaleId: fakeUint256},
					Type: types.Uint256,
				}

				Expect(mappings[salesTopKey]).To(Equal(expectedMetadata))
			})
		})
	})
})
