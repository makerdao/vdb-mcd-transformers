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

package flip_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	mcdStorage "github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/mcd_transformers/transformers/storage/flip"
	"github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
	"github.com/vulcanize/vulcanizedb/libraries/shared/factories/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

var _ = Describe("Flip storage keys loader", func() {
	var (
		storageRepository *test_helpers.MockMakerStorageRepository
		storageKeysLoader storage.KeysLoader
	)

	BeforeEach(func() {
		storageRepository = &test_helpers.MockMakerStorageRepository{}
		storageKeysLoader = flip.NewKeysLoader(storageRepository, constants.GetContractAddress("MCD_FLIP_ETH_A"))
	})

	It("returns value metadata for static keys", func() {
		mappings, err := storageKeysLoader.LoadMappings()

		Expect(err).NotTo(HaveOccurred())
		Expect(mappings[flip.VatKey]).To(Equal(flip.VatMetadata))
		Expect(mappings[flip.IlkKey]).To(Equal(flip.IlkMetadata))
		Expect(mappings[flip.BegKey]).To(Equal(flip.BegMetadata))
		Expect(mappings[flip.TtlAndTauStorageKey]).To(Equal(flip.TtlAndTauMetadata))
		Expect(mappings[flip.KicksKey]).To(Equal(flip.KicksMetadata))
	})

	Describe("bid", func() {
		Describe("when loading bid IDs fails", func() {
			It("returns error", func() {
				storageRepository.GetFlipBidIdsError = fakes.FakeError

				_, err := storageKeysLoader.LoadMappings()

				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(fakes.FakeError))
			})
		})

		Describe("when loading bid IDs succeeds", func() {
			var (
				fakeBidId       = "1"
				fakeHexBidId, _ = shared.ConvertIntStringToHex(fakeBidId)
				bidBidKey       = common.BytesToHash(crypto.Keccak256(common.FromHex(fakeHexBidId + flip.BidsMappingIndex)))
				mappings        map[common.Hash]utils.StorageValueMetadata
			)

			BeforeEach(func() {
				storageRepository.FlipBidIds = []string{fakeBidId}
				var err error
				mappings, err = storageKeysLoader.LoadMappings()
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns value metadata for bid bid", func() {
				expectedMetadata := utils.StorageValueMetadata{
					Name: mcdStorage.BidBid,
					Keys: map[utils.Key]string{constants.BidId: fakeBidId},
					Type: utils.Uint256,
				}

				Expect(mappings[bidBidKey]).To(Equal(expectedMetadata))
			})

			It("returns value metadata for bid lot", func() {
				bidLotKey := utils.GetIncrementedStorageKey(bidBidKey, 1)
				expectedMetadata := utils.StorageValueMetadata{
					Name: mcdStorage.BidLot,
					Keys: map[utils.Key]string{constants.BidId: fakeBidId},
					Type: utils.Uint256,
				}

				Expect(mappings[bidLotKey]).To(Equal(expectedMetadata))
			})

			It("returns value metadata for bid guy + tic + end packed slot", func() {
				bidGuyKey := utils.GetIncrementedStorageKey(bidBidKey, 2)
				expectedMetadata := utils.StorageValueMetadata{
					Name:        mcdStorage.Packed,
					Keys:        map[utils.Key]string{constants.BidId: fakeBidId},
					Type:        utils.PackedSlot,
					PackedTypes: map[int]utils.ValueType{0: utils.Address, 1: utils.Uint48, 2: utils.Uint48},
					PackedNames: map[int]string{0: mcdStorage.BidGuy, 1: mcdStorage.BidTic, 2: mcdStorage.BidEnd},
				}

				Expect(mappings[bidGuyKey]).To(Equal(expectedMetadata))
			})

			It("returns value metadata for bid usr", func() {
				bidUsrKey := utils.GetIncrementedStorageKey(bidBidKey, 3)
				expectedMetadata := utils.StorageValueMetadata{
					Name: mcdStorage.BidUsr,
					Keys: map[utils.Key]string{constants.BidId: fakeBidId},
					Type: utils.Address,
				}

				Expect(mappings[bidUsrKey]).To(Equal(expectedMetadata))
			})

			It("returns value metadata for bid gal", func() {
				bidGalKey := utils.GetIncrementedStorageKey(bidBidKey, 4)
				expectedMetadata := utils.StorageValueMetadata{
					Name: mcdStorage.BidGal,
					Keys: map[utils.Key]string{constants.BidId: fakeBidId},
					Type: utils.Address,
				}

				Expect(mappings[bidGalKey]).To(Equal(expectedMetadata))
			})

			It("returns value metadata for bid tab", func() {
				bidTabKey := utils.GetIncrementedStorageKey(bidBidKey, 5)
				expectedMetadata := utils.StorageValueMetadata{
					Name: mcdStorage.BidTab,
					Keys: map[utils.Key]string{constants.BidId: fakeBidId},
					Type: utils.Uint256,
				}

				Expect(mappings[bidTabKey]).To(Equal(expectedMetadata))
			})
		})
	})
})
