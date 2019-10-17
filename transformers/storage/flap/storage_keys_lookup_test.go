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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	mcdStorage "github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/mcd_transformers/transformers/storage/flap"
	"github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

var _ = Describe("Flap storage mappings", func() {
	var (
		storageRepository *test_helpers.MockMakerStorageRepository
		storageKeysLookup flap.StorageKeysLookup
	)

	BeforeEach(func() {
		storageRepository = &test_helpers.MockMakerStorageRepository{}
		storageKeysLookup = flap.StorageKeysLookup{
			StorageRepository: storageRepository,
			ContractAddress:   test_data.FlapAddress(),
		}
	})

	Describe("looks up static keys", func() {
		It("returns storage value mapping if storage key exists", func() {
			Expect(storageKeysLookup.Lookup(flap.VatStorageKey)).To(Equal(utils.StorageValueMetadata{
				Name: mcdStorage.Vat,
				Keys: nil,
				Type: utils.Address,
			}))
			Expect(storageKeysLookup.Lookup(flap.GemStorageKey)).To(Equal(utils.StorageValueMetadata{
				Name: mcdStorage.Gem,
				Keys: nil,
				Type: utils.Address,
			}))
			Expect(storageKeysLookup.Lookup(flap.BegStorageKey)).To(Equal(utils.StorageValueMetadata{
				Name: mcdStorage.Beg,
				Keys: nil,
				Type: utils.Uint256,
			}))
			Expect(storageKeysLookup.Lookup(flap.TtlAndTauStorageKey)).To(Equal(utils.StorageValueMetadata{
				Name:        mcdStorage.Packed,
				Type:        utils.PackedSlot,
				PackedTypes: map[int]utils.ValueType{0: utils.Uint48, 1: utils.Uint48},
				PackedNames: map[int]string{0: mcdStorage.Ttl, 1: mcdStorage.Tau},
			}))
			Expect(storageKeysLookup.Lookup(flap.KicksStorageKey)).To(Equal(utils.StorageValueMetadata{
				Name: mcdStorage.Kicks,
				Keys: nil,
				Type: utils.Uint256,
			}))
			Expect(storageKeysLookup.Lookup(flap.LiveStorageKey)).To(Equal(utils.StorageValueMetadata{
				Name: mcdStorage.Live,
				Keys: nil,
				Type: utils.Uint256,
			}))
		})

		It("returns storage value mapping if keccak hashed storage key exists", func() {
			Expect(storageKeysLookup.Lookup(crypto.Keccak256Hash(flap.VatStorageKey[:]))).To(Equal(utils.StorageValueMetadata{
				Name: mcdStorage.Vat,
				Keys: nil,
				Type: utils.Address,
			}))
			Expect(storageKeysLookup.Lookup(crypto.Keccak256Hash(flap.GemStorageKey[:]))).To(Equal(utils.StorageValueMetadata{
				Name: mcdStorage.Gem,
				Keys: nil,
				Type: utils.Address,
			}))
			Expect(storageKeysLookup.Lookup(crypto.Keccak256Hash(flap.BegStorageKey[:]))).To(Equal(utils.StorageValueMetadata{
				Name: mcdStorage.Beg,
				Keys: nil,
				Type: utils.Uint256,
			}))
			Expect(storageKeysLookup.Lookup(crypto.Keccak256Hash(flap.TtlAndTauStorageKey[:]))).To(Equal(utils.StorageValueMetadata{
				Name:        mcdStorage.Packed,
				Type:        utils.PackedSlot,
				PackedTypes: map[int]utils.ValueType{0: utils.Uint48, 1: utils.Uint48},
				PackedNames: map[int]string{0: mcdStorage.Ttl, 1: mcdStorage.Tau},
			}))
			Expect(storageKeysLookup.Lookup(crypto.Keccak256Hash(flap.KicksStorageKey[:]))).To(Equal(utils.StorageValueMetadata{
				Name: mcdStorage.Kicks,
				Keys: nil,
				Type: utils.Uint256,
			}))
			Expect(storageKeysLookup.Lookup(crypto.Keccak256Hash(flap.LiveStorageKey[:]))).To(Equal(utils.StorageValueMetadata{
				Name: mcdStorage.Live,
				Keys: nil,
				Type: utils.Uint256,
			}))
		})

		It("returns an error if the key doesn't exist", func() {
			emptyMetadata, err := storageKeysLookup.Lookup(fakes.FakeHash)

			Expect(emptyMetadata).To(Equal(utils.StorageValueMetadata{}))
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrStorageKeyNotFound{Key: fakes.FakeHash.Hex()}))
		})
	})

	Describe("looks up dynamic keys", func() {
		It("returns error if getting bid ids fails", func() {
			storageRepository.GetFlapBidIdsError = fakes.FakeError

			_, err := storageKeysLookup.Lookup(fakes.FakeHash)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})

		It("refreshes the bid keys if the given key is not found", func() {
			storageKeysLookup.Lookup(fakes.FakeHash)

			Expect(storageRepository.GetFlapBidIdsCalled).To(BeTrue())
		})

		Describe("flap bids", func() {
			var (
				bidId         = strconv.FormatInt(rand.Int63(), 10)
				bidIdHex, _   = shared.ConvertIntStringToHex(bidId)
				flapBidBidKey = common.BytesToHash(
					crypto.Keccak256(
						common.FromHex(bidIdHex + flap.BidsIndex)))
			)

			BeforeEach(func() {
				storageRepository.FlapBidIds = []string{bidId}
			})

			It("gets bid metadata", func() {
				metadata, err := storageKeysLookup.Lookup(flapBidBidKey)
				Expect(err).NotTo(HaveOccurred())

				Expect(metadata).To(Equal(utils.StorageValueMetadata{
					Name: mcdStorage.BidBid,
					Keys: map[utils.Key]string{constants.BidId: bidId},
					Type: utils.Uint256,
				}))
			})

			It("gets lot metadata", func() {
				flapBidLotKey := storage.GetIncrementedKey(flapBidBidKey, 1)
				metadata, err := storageKeysLookup.Lookup(flapBidLotKey)
				Expect(err).NotTo(HaveOccurred())

				Expect(metadata).To(Equal(utils.StorageValueMetadata{
					Name: mcdStorage.BidLot,
					Keys: map[utils.Key]string{constants.BidId: bidId},
					Type: utils.Uint256,
				}))
			})

			It("returns value metadata for bid guy + tic + end packed slot", func() {
				bidGuyKey := storage.GetIncrementedKey(flapBidBidKey, 2)
				expectedMetadata := utils.StorageValueMetadata{
					Name:        mcdStorage.Packed,
					Keys:        map[utils.Key]string{constants.BidId: bidId},
					Type:        utils.PackedSlot,
					PackedTypes: map[int]utils.ValueType{0: utils.Address, 1: utils.Uint48, 2: utils.Uint48},
					PackedNames: map[int]string{0: mcdStorage.BidGuy, 1: mcdStorage.BidTic, 2: mcdStorage.BidEnd},
				}
				Expect(storageKeysLookup.Lookup(bidGuyKey)).To(Equal(expectedMetadata))
			})
		})
	})
})
