package dog_test

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/dog"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
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
