package medianizer_test

import (
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/medianizer"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
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

		Expect(mappings[medianizer.ValAndAgeStorageKey]).To(Equal(types.ValueMetadata{
			Name:        mcdStorage.Packed,
			Type:        types.PackedSlot,
			PackedTypes: map[int]types.ValueType{0: types.Uint128, 1: types.Uint48},
			PackedNames: map[int]string{0: medianizer.Val, 1: medianizer.Age},
		}))

		Expect(mappings[medianizer.BarKey]).To(Equal(medianizer.BarMetadata))
	})
})
