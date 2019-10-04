package cat_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage/cat"
	"github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
)

var _ = Describe("Cat storage mappings", func() {

	var (
		storageRepository *test_helpers.MockMakerStorageRepository
		mappings          cat.CatMappings
	)

	BeforeEach(func() {
		storageRepository = &test_helpers.MockMakerStorageRepository{}
		mappings = cat.CatMappings{StorageRepository: storageRepository}
	})

	Describe("looking up static keys", func() {
		It("returns value metadata if key exists", func() {
			Expect(mappings.Lookup(cat.LiveKey)).To(Equal(cat.LiveMetadata))
			Expect(mappings.Lookup(cat.VatKey)).To(Equal(cat.VatMetadata))
			Expect(mappings.Lookup(cat.VowKey)).To(Equal(cat.VowMetadata))
		})

		It("returns value metadata for keccak hashed storage keys", func() {
			Expect(mappings.Lookup(crypto.Keccak256Hash(cat.LiveKey[:]))).To(Equal(cat.LiveMetadata))
			Expect(mappings.Lookup(crypto.Keccak256Hash(cat.VatKey[:]))).To(Equal(cat.VatMetadata))
			Expect(mappings.Lookup(crypto.Keccak256Hash(cat.VowKey[:]))).To(Equal(cat.VowMetadata))
		})

		It("returns error if key does not exist", func() {
			_, err := mappings.Lookup(fakes.FakeHash)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrStorageKeyNotFound{Key: fakes.FakeHash.Hex()}))
		})
	})

	Describe("looking up dynamic keys", func() {
		It("refreshes mappings from repository if key not found", func() {
			_, _ = mappings.Lookup(fakes.FakeHash)

			Expect(storageRepository.GetIlksCalled).To(BeTrue())
		})

		It("returns error if ilks lookup fails", func() {
			storageRepository.GetIlksError = fakes.FakeError

			_, err := mappings.Lookup(fakes.FakeHash)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})

		Describe("ilk", func() {
			var ilkFlipKey = common.BytesToHash(crypto.Keccak256(common.FromHex(test_helpers.FakeIlk + cat.IlksMappingIndex)))

			BeforeEach(func() {
				storageRepository.Ilks = []string{test_helpers.FakeIlk}
			})

			It("returns value metadata for ilk flip", func() {
				expectedMetadata := utils.StorageValueMetadata{
					Name: cat.IlkFlip,
					Keys: map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk},
					Type: utils.Address,
				}
				Expect(mappings.Lookup(ilkFlipKey)).To(Equal(expectedMetadata))
			})

			It("returns value metadata for ilk chop", func() {
				ilkChopKey := storage.GetIncrementedKey(ilkFlipKey, 1)
				expectedMetadata := utils.StorageValueMetadata{
					Name: cat.IlkChop,
					Keys: map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk},
					Type: utils.Uint256,
				}
				Expect(mappings.Lookup(ilkChopKey)).To(Equal(expectedMetadata))
			})

			It("returns value metadata for ilk lump", func() {
				ilkLumpKey := storage.GetIncrementedKey(ilkFlipKey, 2)
				expectedMetadata := utils.StorageValueMetadata{
					Name: cat.IlkLump,
					Keys: map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk},
					Type: utils.Uint256,
				}
				Expect(mappings.Lookup(ilkLumpKey)).To(Equal(expectedMetadata))
			})
		})
	})
})
