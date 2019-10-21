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
	mcdStorage "github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/mcd_transformers/transformers/storage/cdp_manager"
	"github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/storage/utilities"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/rand"
	"strconv"
)

var _ = Describe("CDP Manager storage keys loader", func() {
	var (
		storageRepository *test_helpers.MockMakerStorageRepository
		storageKeysLoader mcdStorage.KeysLoader
	)

	BeforeEach(func() {
		storageRepository = &test_helpers.MockMakerStorageRepository{}
		storageKeysLoader = cdp_manager.NewKeysLoader(storageRepository)
	})

	It("returns value metadata for static keys", func() {
		mappings, err := storageKeysLoader.LoadMappings()

		Expect(err).NotTo(HaveOccurred())
		Expect(mappings[cdp_manager.VatKey]).To(Equal(cdp_manager.VatMetadata))
		Expect(mappings[cdp_manager.CdpiKey]).To(Equal(cdp_manager.CdpiMetadata))
	})

	Describe("looking up dynamic keys", func() {
		Describe("Mappings with CDPI as key", func() {
			Describe("when getting CDPIs fails", func() {
				It("returns error", func() {
					storageRepository.GetCdpisError = fakes.FakeError

					_, err := storageKeysLoader.LoadMappings()

					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(fakes.FakeError))
				})
			})

			Describe("when getting CDPIs succeeds", func() {

				var (
					cdpi       = strconv.FormatInt(rand.Int63(), 10)
					cdpiHex, _ = shared.ConvertIntStringToHex(cdpi)
					mappings   map[common.Hash]utils.StorageValueMetadata
				)

				BeforeEach(func() {
					storageRepository.Cdpis = []string{cdpi}
					var err error
					mappings, err = storageKeysLoader.LoadMappings()
					Expect(err).NotTo(HaveOccurred())
				})

				It("gets Urns metadata", func() {
					urnsKey := common.BytesToHash(crypto.Keccak256(common.FromHex(cdpiHex + cdp_manager.UrnsMappingIndex)))

					Expect(mappings[urnsKey]).To(Equal(utils.StorageValueMetadata{
						Name: cdp_manager.Urns,
						Keys: map[utils.Key]string{constants.Cdpi: cdpi},
						Type: utils.Address,
					}))
				})

				Describe("List mappings", func() {
					listPrevKey := common.BytesToHash(
						crypto.Keccak256(common.FromHex(cdpiHex + cdp_manager.ListMappingIndex)))

					It("gets prev metadata", func() {
						Expect(mappings[listPrevKey]).To(Equal(utils.StorageValueMetadata{
							Name: cdp_manager.ListPrev,
							Keys: map[utils.Key]string{constants.Cdpi: cdpi},
							Type: utils.Uint256,
						}))
					})

					It("gets next metadata", func() {
						listNextKey := storage.GetIncrementedKey(listPrevKey, 1)

						Expect(mappings[listNextKey]).To(Equal(utils.StorageValueMetadata{
							Name: cdp_manager.ListNext,
							Keys: map[utils.Key]string{constants.Cdpi: cdpi},
							Type: utils.Uint256,
						}))
					})
				})

				It("gets Owner metadata", func() {
					ownsKey := common.BytesToHash(crypto.Keccak256(common.FromHex(cdpiHex + cdp_manager.OwnsMappingIndex)))

					Expect(mappings[ownsKey]).To(Equal(utils.StorageValueMetadata{
						Name: cdp_manager.Owns,
						Keys: map[utils.Key]string{constants.Cdpi: cdpi},
						Type: utils.Address,
					}))
				})

				It("gets Ilks metadata", func() {
					ilksKey := common.BytesToHash(crypto.Keccak256(common.FromHex(cdpiHex + cdp_manager.IlksMappingIndex)))

					Expect(mappings[ilksKey]).To(Equal(utils.StorageValueMetadata{
						Name: cdp_manager.Ilks,
						Keys: map[utils.Key]string{constants.Cdpi: cdpi},
						Type: utils.Bytes32,
					}))
				})
			})
		})

		Describe("Mappings with Owner as key", func() {
			Describe("when getting owners fails", func() {
				It("returns error", func() {
					storageRepository.GetOwnersError = fakes.FakeError

					_, err := storageKeysLoader.LoadMappings()

					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(fakes.FakeError))
				})
			})

			Describe("when getting owners succeeds", func() {
				var (
					owns          = test_helpers.FakeAddress
					paddedOwns, _ = utilities.PadAddress(owns)
					mappings      map[common.Hash]utils.StorageValueMetadata
				)

				BeforeEach(func() {
					storageRepository.Owners = []string{owns}
					var err error
					mappings, err = storageKeysLoader.LoadMappings()
					Expect(err).NotTo(HaveOccurred())
				})

				It("gets First metadata", func() {
					firstKey := common.BytesToHash(crypto.Keccak256(common.FromHex(paddedOwns + cdp_manager.FirstMappingIndex)))

					Expect(mappings[firstKey]).To(Equal(utils.StorageValueMetadata{
						Name: cdp_manager.First,
						Keys: map[utils.Key]string{constants.Owner: owns},
						Type: utils.Uint256,
					}))
				})

				It("gets Last metadata", func() {
					lastKey := common.BytesToHash(crypto.Keccak256(common.FromHex(paddedOwns + cdp_manager.LastMappingIndex)))

					Expect(mappings[lastKey]).To(Equal(utils.StorageValueMetadata{
						Name: cdp_manager.Last,
						Keys: map[utils.Key]string{constants.Owner: owns},
						Type: utils.Uint256,
					}))
				})

				It("gets Count metadata", func() {
					countKey := common.BytesToHash(crypto.Keccak256(common.FromHex(paddedOwns + cdp_manager.CountMappingIndex)))

					Expect(mappings[countKey]).To(Equal(utils.StorageValueMetadata{
						Name: cdp_manager.Count,
						Keys: map[utils.Key]string{constants.Owner: owns},
						Type: utils.Uint256,
					}))
				})
			})
		})
	})
})
