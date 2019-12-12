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

package flap_test

import (
	"math/rand"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/flap"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	vdbStorage "github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Flap storage keys loader", func() {
	var (
		storageRepository *test_helpers.MockMakerStorageRepository
		storageKeysLoader storage.KeysLoader
	)

	BeforeEach(func() {
		storageRepository = &test_helpers.MockMakerStorageRepository{}
		storageKeysLoader = flap.NewKeysLoader(storageRepository, test_data.FlapAddress())
	})

	It("returns storage value metadata for static keys", func() {
		mappings, err := storageKeysLoader.LoadMappings()

		Expect(err).NotTo(HaveOccurred())
		Expect(mappings[flap.VatStorageKey]).To(Equal(vdbStorage.StorageValueMetadata{
			Name: mcdStorage.Vat,
			Keys: nil,
			Type: vdbStorage.Address,
		}))
		Expect(mappings[flap.GemStorageKey]).To(Equal(vdbStorage.StorageValueMetadata{
			Name: mcdStorage.Gem,
			Keys: nil,
			Type: vdbStorage.Address,
		}))
		Expect(mappings[flap.BegStorageKey]).To(Equal(vdbStorage.StorageValueMetadata{
			Name: mcdStorage.Beg,
			Keys: nil,
			Type: vdbStorage.Uint256,
		}))
		Expect(mappings[flap.TtlAndTauStorageKey]).To(Equal(vdbStorage.StorageValueMetadata{
			Name:        mcdStorage.Packed,
			Type:        vdbStorage.PackedSlot,
			PackedTypes: map[int]vdbStorage.ValueType{0: vdbStorage.Uint48, 1: vdbStorage.Uint48},
			PackedNames: map[int]string{0: mcdStorage.Ttl, 1: mcdStorage.Tau},
		}))
		Expect(mappings[flap.KicksStorageKey]).To(Equal(vdbStorage.StorageValueMetadata{
			Name: mcdStorage.Kicks,
			Keys: nil,
			Type: vdbStorage.Uint256,
		}))
		Expect(mappings[flap.LiveStorageKey]).To(Equal(vdbStorage.StorageValueMetadata{
			Name: mcdStorage.Live,
			Keys: nil,
			Type: vdbStorage.Uint256,
		}))
	})

	Describe("flap bids", func() {
		Describe("when getting flap bid IDs fails", func() {
			It("returns error", func() {
				storageRepository.GetFlapBidIdsError = fakes.FakeError

				_, err := storageKeysLoader.LoadMappings()

				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(fakes.FakeError))
			})
		})

		Describe("when getting flap bid IDs succeeds", func() {
			var (
				bidId         = strconv.FormatInt(rand.Int63(), 10)
				bidIdHex, _   = shared.ConvertIntStringToHex(bidId)
				flapBidBidKey = common.BytesToHash(
					crypto.Keccak256(
						common.FromHex(bidIdHex + flap.BidsIndex)))
				mappings map[common.Hash]vdbStorage.StorageValueMetadata
			)

			BeforeEach(func() {
				storageRepository.FlapBidIds = []string{bidId}
				var err error
				mappings, err = storageKeysLoader.LoadMappings()
				Expect(err).NotTo(HaveOccurred())
			})

			It("gets bid metadata", func() {
				expectedMetadata := vdbStorage.StorageValueMetadata{
					Name: mcdStorage.BidBid,
					Keys: map[vdbStorage.Key]string{constants.BidId: bidId},
					Type: vdbStorage.Uint256,
				}

				Expect(mappings[flapBidBidKey]).To(Equal(expectedMetadata))
			})

			It("gets lot metadata", func() {
				flapBidLotKey := vdbStorage.GetIncrementedStorageKey(flapBidBidKey, 1)
				expectedMetadata := vdbStorage.StorageValueMetadata{
					Name: mcdStorage.BidLot,
					Keys: map[vdbStorage.Key]string{constants.BidId: bidId},
					Type: vdbStorage.Uint256,
				}

				Expect(mappings[flapBidLotKey]).To(Equal(expectedMetadata))
			})

			It("returns value metadata for bid guy + tic + end packed slot", func() {
				bidGuyKey := vdbStorage.GetIncrementedStorageKey(flapBidBidKey, 2)
				expectedMetadata := vdbStorage.StorageValueMetadata{
					Name:        mcdStorage.Packed,
					Keys:        map[vdbStorage.Key]string{constants.BidId: bidId},
					Type:        vdbStorage.PackedSlot,
					PackedTypes: map[int]vdbStorage.ValueType{0: vdbStorage.Address, 1: vdbStorage.Uint48, 2: vdbStorage.Uint48},
					PackedNames: map[int]string{0: mcdStorage.BidGuy, 1: mcdStorage.BidTic, 2: mcdStorage.BidEnd},
				}

				Expect(mappings[bidGuyKey]).To(Equal(expectedMetadata))
			})
		})
	})
})
