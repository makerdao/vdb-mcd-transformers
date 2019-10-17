// VulcanizeDB
// Copyright © 2018 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package cat_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage/cat"
	"github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

var _ = Describe("Cat storage mappings", func() {

	var (
		storageRepository *test_helpers.MockMakerStorageRepository
		storageKeysLookup cat.StorageKeysLookup
	)

	BeforeEach(func() {
		storageRepository = &test_helpers.MockMakerStorageRepository{}
		storageKeysLookup = cat.StorageKeysLookup{StorageRepository: storageRepository}
	})

	Describe("looking up static keys", func() {
		It("returns value metadata if key exists", func() {
			Expect(storageKeysLookup.Lookup(cat.LiveKey)).To(Equal(cat.LiveMetadata))
			Expect(storageKeysLookup.Lookup(cat.VatKey)).To(Equal(cat.VatMetadata))
			Expect(storageKeysLookup.Lookup(cat.VowKey)).To(Equal(cat.VowMetadata))
		})

		It("returns value metadata for keccak hashed storage keys", func() {
			Expect(storageKeysLookup.Lookup(crypto.Keccak256Hash(cat.LiveKey[:]))).To(Equal(cat.LiveMetadata))
			Expect(storageKeysLookup.Lookup(crypto.Keccak256Hash(cat.VatKey[:]))).To(Equal(cat.VatMetadata))
			Expect(storageKeysLookup.Lookup(crypto.Keccak256Hash(cat.VowKey[:]))).To(Equal(cat.VowMetadata))
		})

		It("returns error if key does not exist", func() {
			_, err := storageKeysLookup.Lookup(fakes.FakeHash)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrStorageKeyNotFound{Key: fakes.FakeHash.Hex()}))
		})
	})

	Describe("looking up dynamic keys", func() {
		It("refreshes mappings from repository if key not found", func() {
			_, _ = storageKeysLookup.Lookup(fakes.FakeHash)

			Expect(storageRepository.GetIlksCalled).To(BeTrue())
		})

		It("returns error if ilks lookup fails", func() {
			storageRepository.GetIlksError = fakes.FakeError

			_, err := storageKeysLookup.Lookup(fakes.FakeHash)

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
				Expect(storageKeysLookup.Lookup(ilkFlipKey)).To(Equal(expectedMetadata))
			})

			It("returns value metadata for ilk chop", func() {
				ilkChopKey := storage.GetIncrementedKey(ilkFlipKey, 1)
				expectedMetadata := utils.StorageValueMetadata{
					Name: cat.IlkChop,
					Keys: map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk},
					Type: utils.Uint256,
				}
				Expect(storageKeysLookup.Lookup(ilkChopKey)).To(Equal(expectedMetadata))
			})

			It("returns value metadata for ilk lump", func() {
				ilkLumpKey := storage.GetIncrementedKey(ilkFlipKey, 2)
				expectedMetadata := utils.StorageValueMetadata{
					Name: cat.IlkLump,
					Keys: map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk},
					Type: utils.Uint256,
				}
				Expect(storageKeysLookup.Lookup(ilkLumpKey)).To(Equal(expectedMetadata))
			})
		})
	})
})
