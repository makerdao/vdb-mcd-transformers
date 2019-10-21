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
	mcdStorage "github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/mcd_transformers/transformers/storage/cat"
	"github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

var _ = Describe("Cat storage keys loader", func() {
	var (
		storageRepository *test_helpers.MockMakerStorageRepository
		storageKeysLoader mcdStorage.KeysLoader
	)

	BeforeEach(func() {
		storageRepository = &test_helpers.MockMakerStorageRepository{}
		storageKeysLoader = cat.NewKeysLoader(storageRepository)
	})

	It("returns value metadata for static keys", func() {
		mappings, err := storageKeysLoader.LoadMappings()

		Expect(err).NotTo(HaveOccurred())
		Expect(mappings[cat.LiveKey]).To(Equal(cat.LiveMetadata))
		Expect(mappings[cat.VatKey]).To(Equal(cat.VatMetadata))
		Expect(mappings[cat.VowKey]).To(Equal(cat.VowMetadata))
	})

	Describe("ilk", func() {
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
				ilkFlipKey = common.BytesToHash(crypto.Keccak256(common.FromHex(test_helpers.FakeIlk + cat.IlksMappingIndex)))
				mappings   map[common.Hash]utils.StorageValueMetadata
			)

			BeforeEach(func() {
				storageRepository.Ilks = []string{test_helpers.FakeIlk}
				var err error
				mappings, err = storageKeysLoader.LoadMappings()
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns value metadata for ilk flip", func() {
				expectedMetadata := utils.StorageValueMetadata{
					Name: cat.IlkFlip,
					Keys: map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk},
					Type: utils.Address,
				}

				Expect(mappings[ilkFlipKey]).To(Equal(expectedMetadata))
			})

			It("returns value metadata for ilk chop", func() {
				ilkChopKey := storage.GetIncrementedKey(ilkFlipKey, 1)
				expectedMetadata := utils.StorageValueMetadata{
					Name: cat.IlkChop,
					Keys: map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk},
					Type: utils.Uint256,
				}

				Expect(mappings[ilkChopKey]).To(Equal(expectedMetadata))
			})

			It("returns value metadata for ilk lump", func() {
				ilkLumpKey := storage.GetIncrementedKey(ilkFlipKey, 2)
				expectedMetadata := utils.StorageValueMetadata{
					Name: cat.IlkLump,
					Keys: map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk},
					Type: utils.Uint256,
				}

				Expect(mappings[ilkLumpKey]).To(Equal(expectedMetadata))
			})
		})
	})
})
