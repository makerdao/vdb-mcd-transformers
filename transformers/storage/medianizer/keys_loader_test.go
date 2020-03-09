package medianizer_test

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/medianizer"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Medianizer storage keys loader", func() {
	var (
		storageRepository *test_helpers.MockMakerStorageRepository
		storageKeysLoader storage.KeysLoader
	)
	BeforeEach(func() {
		storageRepository = &test_helpers.MockMakerStorageRepository{}
		storageKeysLoader = medianizer.NewKeysLoader(storageRepository, test_data.MedianizerAddress())
	})

	It("returns value metadata for static keys", func() {
		mappings, err := storageKeysLoader.LoadMappings()

		Expect(err).NotTo(HaveOccurred())
		Expect(mappings[medianizer.ValKey]).To(Equal(medianizer.ValMetadata))
		Expect(mappings[medianizer.AgeKey]).To(Equal(medianizer.AgeMetadata))
		Expect(mappings[medianizer.BarKey]).To(Equal(medianizer.BarMetadata))
	})
})
