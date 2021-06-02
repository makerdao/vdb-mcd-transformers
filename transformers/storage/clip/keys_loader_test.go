package clip_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities/wards"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
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
		Expect(mappings[clip.ChipAndTipStorageKey]).To(Equal(clip.ChipAndTipMetadata))
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
	It("returns value for dynamic size array keys", func() {
		mappings, err := storageKeysLoader.LoadMappings()

		Expect(err).NotTo(HaveOccurred())
		Expect(mappings[clip.ActiveKey]).To(Equal(clip.ActiveMetadata))
	})
})
