package clip_test

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
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
		storageKeysLoader = clip.NewKeysLoader(storageRepository, test_data.Clip130Address())
	})

	It("returns value metadata for static keys", func() {
		mappings, err := storageKeysLoader.LoadMappings()

		Expect(err).NotTo(HaveOccurred())
		Expect(mappings[clip.IlkKey]).To(Equal(clip.IlkMetadata))
		Expect(mappings[clip.VatKey]).To(Equal(clip.VatMetadata))
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
})
