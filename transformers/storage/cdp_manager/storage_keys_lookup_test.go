// VulcanizeDB
// Copyright © 2019 Vulcanize

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

package cdp_manager_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage/cdp_manager"
	"github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/storage/utilities"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/rand"
	"strconv"
)

var _ = Describe("CDP Manager storage mappings", func() {
	var (
		storageRepository *test_helpers.MockMakerStorageRepository
		storageKeysLookup cdp_manager.StorageKeysLookup
	)

	BeforeEach(func() {
		storageRepository = &test_helpers.MockMakerStorageRepository{}
		storageKeysLookup = cdp_manager.StorageKeysLookup{StorageRepository: storageRepository}
	})

	Describe("looking up static keys", func() {
		It("returns value metadata if key exists", func() {
			Expect(storageKeysLookup.Lookup(cdp_manager.VatKey)).To(Equal(cdp_manager.VatMetadata))
			Expect(storageKeysLookup.Lookup(cdp_manager.CdpiKey)).To(Equal(cdp_manager.CdpiMetadata))
		})

		It("returns error if key does not exist", func() {
			_, err := storageKeysLookup.Lookup(common.HexToHash(fakes.FakeHash.Hex()))

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrStorageKeyNotFound{Key: fakes.FakeHash.Hex()}))
		})
	})

	Describe("looking up dynamic keys", func() {
		Describe("Mappings with CDPI as key", func() {
			var (
				cdpi       = strconv.FormatInt(rand.Int63(), 10)
				cdpiHex, _ = shared.ConvertIntStringToHex(cdpi)
			)

			BeforeEach(func() {
				storageRepository.Cdpis = []string{cdpi}
			})

			It("gets Urns metadata", func() {
				urnsKey := common.BytesToHash(crypto.Keccak256(common.FromHex(cdpiHex + cdp_manager.UrnsMappingIndex)))
				metadata, err := storageKeysLookup.Lookup(urnsKey)
				Expect(err).NotTo(HaveOccurred())

				Expect(metadata).To(Equal(utils.StorageValueMetadata{
					Name: cdp_manager.CdpManagerUrns,
					Keys: map[utils.Key]string{constants.Cdpi: cdpi},
					Type: utils.Address,
				}))
			})

			Describe("List mappings", func() {
				listPrevKey := common.BytesToHash(
					crypto.Keccak256(common.FromHex(cdpiHex + cdp_manager.ListMappingIndex)))

				It("gets prev metadata", func() {
					metadata, err := storageKeysLookup.Lookup(listPrevKey)
					Expect(err).NotTo(HaveOccurred())

					Expect(metadata).To(Equal(utils.StorageValueMetadata{
						Name: cdp_manager.CdpManagerListPrev,
						Keys: map[utils.Key]string{constants.Cdpi: cdpi},
						Type: utils.Uint256,
					}))
				})

				It("gets next metadata", func() {
					listNextKey := storage.GetIncrementedKey(listPrevKey, 1)
					metadata, err := storageKeysLookup.Lookup(listNextKey)
					Expect(err).NotTo(HaveOccurred())

					Expect(metadata).To(Equal(utils.StorageValueMetadata{
						Name: cdp_manager.CdpManagerListNext,
						Keys: map[utils.Key]string{constants.Cdpi: cdpi},
						Type: utils.Uint256,
					}))
				})
			})

			It("gets Owner metadata", func() {
				ownsKey := common.BytesToHash(crypto.Keccak256(common.FromHex(cdpiHex + cdp_manager.OwnsMappingIndex)))
				metadata, err := storageKeysLookup.Lookup(ownsKey)
				Expect(err).NotTo(HaveOccurred())

				Expect(metadata).To(Equal(utils.StorageValueMetadata{
					Name: cdp_manager.CdpManagerOwns,
					Keys: map[utils.Key]string{constants.Cdpi: cdpi},
					Type: utils.Address,
				}))
			})

			It("gets Ilks metadata", func() {
				ilksKey := common.BytesToHash(crypto.Keccak256(common.FromHex(cdpiHex + cdp_manager.IlksMappingIndex)))
				metadata, err := storageKeysLookup.Lookup(ilksKey)
				Expect(err).NotTo(HaveOccurred())

				Expect(metadata).To(Equal(utils.StorageValueMetadata{
					Name: cdp_manager.CdpManagerIlks,
					Keys: map[utils.Key]string{constants.Cdpi: cdpi},
					Type: utils.Bytes32,
				}))
			})
		})

		Describe("Mappings with Owner as key", func() {
			var (
				owns          = test_helpers.FakeAddress
				paddedOwns, _ = utilities.PadAddress(owns)
			)

			BeforeEach(func() {
				storageRepository.Owners = []string{owns}
			})

			It("gets First metadata", func() {
				firstKey := common.BytesToHash(crypto.Keccak256(common.FromHex(paddedOwns + cdp_manager.FirstMappingIndex)))
				metadata, err := storageKeysLookup.Lookup(firstKey)
				Expect(err).NotTo(HaveOccurred())

				Expect(metadata).To(Equal(utils.StorageValueMetadata{
					Name: cdp_manager.CdpManagerFirst,
					Keys: map[utils.Key]string{constants.Owner: owns},
					Type: utils.Uint256,
				}))
			})

			It("gets Last metadata", func() {
				lastKey := common.BytesToHash(crypto.Keccak256(common.FromHex(paddedOwns + cdp_manager.LastMappingIndex)))
				metadata, err := storageKeysLookup.Lookup(lastKey)
				Expect(err).NotTo(HaveOccurred())

				Expect(metadata).To(Equal(utils.StorageValueMetadata{
					Name: cdp_manager.CdpManagerLast,
					Keys: map[utils.Key]string{constants.Owner: owns},
					Type: utils.Uint256,
				}))
			})

			It("gets Count metadata", func() {
				countKey := common.BytesToHash(crypto.Keccak256(common.FromHex(paddedOwns + cdp_manager.CountMappingIndex)))
				metadata, err := storageKeysLookup.Lookup(countKey)
				Expect(err).NotTo(HaveOccurred())

				Expect(metadata).To(Equal(utils.StorageValueMetadata{
					Name: cdp_manager.CdpManagerCount,
					Keys: map[utils.Key]string{constants.Owner: owns},
					Type: utils.Uint256,
				}))
			})
		})
	})
})
