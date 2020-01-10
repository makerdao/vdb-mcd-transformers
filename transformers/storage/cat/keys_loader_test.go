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
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/cat"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	vdbStorage "github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cat storage keys loader", func() {
	var (
		storageRepository *test_helpers.MockMakerStorageRepository
		storageKeysLoader storage.KeysLoader
	)

	BeforeEach(func() {
		storageRepository = &test_helpers.MockMakerStorageRepository{}
		storageKeysLoader = cat.NewKeysLoader(storageRepository, test_data.CatAddress())
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
				mappings   map[common.Hash]vdbStorage.ValueMetadata
			)

			BeforeEach(func() {
				storageRepository.Ilks = []string{test_helpers.FakeIlk}
				var err error
				mappings, err = storageKeysLoader.LoadMappings()
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns value metadata for ilk flip", func() {
				expectedMetadata := vdbStorage.ValueMetadata{
					Name: cat.IlkFlip,
					Keys: map[vdbStorage.Key]string{constants.Ilk: test_helpers.FakeIlk},
					Type: vdbStorage.Address,
				}

				Expect(mappings[ilkFlipKey]).To(Equal(expectedMetadata))
			})

			It("returns value metadata for ilk chop", func() {
				ilkChopKey := vdbStorage.GetIncrementedKey(ilkFlipKey, 1)
				expectedMetadata := vdbStorage.ValueMetadata{
					Name: cat.IlkChop,
					Keys: map[vdbStorage.Key]string{constants.Ilk: test_helpers.FakeIlk},
					Type: vdbStorage.Uint256,
				}

				Expect(mappings[ilkChopKey]).To(Equal(expectedMetadata))
			})

			It("returns value metadata for ilk lump", func() {
				ilkLumpKey := vdbStorage.GetIncrementedKey(ilkFlipKey, 2)
				expectedMetadata := vdbStorage.ValueMetadata{
					Name: cat.IlkLump,
					Keys: map[vdbStorage.Key]string{constants.Ilk: test_helpers.FakeIlk},
					Type: vdbStorage.Uint256,
				}

				Expect(mappings[ilkLumpKey]).To(Equal(expectedMetadata))
			})
		})
	})

	Describe("wards", func() {
		It("returns value metadata for wards", func() {
			wardsUser := fakes.FakeAddress.Hex()
			storageRepository.WardsKeys = []string{wardsUser}
			paddedWardsUser := "0x000000000000000000000000" + wardsUser[2:]
			wardsKey := common.BytesToHash(crypto.Keccak256(common.FromHex(paddedWardsUser + cat.WardsMappingIndex)))
			expectedMetadata := vdbStorage.ValueMetadata{
				Name: cat.Wards,
				Keys: map[vdbStorage.Key]string{constants.User: fakes.FakeAddress.Hex()},
				Type: vdbStorage.Uint256,
			}

			mappings, err := storageKeysLoader.LoadMappings()
			Expect(err).NotTo(HaveOccurred())

			Expect(storageRepository.GetWardsKeysCalledWith).To(Equal(test_data.CatAddress()))
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
})
