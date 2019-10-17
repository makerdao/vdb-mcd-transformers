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

package flop_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	mcdStorage "github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/mcd_transformers/transformers/storage/flop"
	"github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

var _ = Describe("Flop storage mappings", func() {
	var (
		storageRepository *test_helpers.MockMakerStorageRepository
		storageKeysLookup flop.StorageKeysLookup
	)

	BeforeEach(func() {
		storageRepository = &test_helpers.MockMakerStorageRepository{}
		storageKeysLookup = flop.StorageKeysLookup{StorageRepository: storageRepository, ContractAddress: "0x668001c75a9c02d6b10c7a17dbd8aa4afff95037"}
	})

	Describe("looking up static keys", func() {
		It("returns value metadata if key exists", func() {
			Expect(storageKeysLookup.Lookup(flop.VatKey)).To(Equal(flop.VatMetadata))
			Expect(storageKeysLookup.Lookup(flop.GemKey)).To(Equal(flop.GemMetadata))
			Expect(storageKeysLookup.Lookup(flop.BegKey)).To(Equal(flop.BegMetadata))
			Expect(storageKeysLookup.Lookup(flop.PadKey)).To(Equal(flop.PadMetadata))
			Expect(storageKeysLookup.Lookup(flop.TtlAndTauKey)).To(Equal(flop.TtlAndTauMetadata))
			Expect(storageKeysLookup.Lookup(flop.KicksKey)).To(Equal(flop.KicksMetadata))
			Expect(storageKeysLookup.Lookup(flop.LiveKey)).To(Equal(flop.LiveMetadata))
		})

		It("returns value metadata if keccak hash of key exists", func() {
			Expect(storageKeysLookup.Lookup(crypto.Keccak256Hash(flop.VatKey[:]))).To(Equal(flop.VatMetadata))
			Expect(storageKeysLookup.Lookup(crypto.Keccak256Hash(flop.GemKey[:]))).To(Equal(flop.GemMetadata))
			Expect(storageKeysLookup.Lookup(crypto.Keccak256Hash(flop.BegKey[:]))).To(Equal(flop.BegMetadata))
			Expect(storageKeysLookup.Lookup(crypto.Keccak256Hash(flop.TtlAndTauKey[:]))).To(Equal(flop.TtlAndTauMetadata))
			Expect(storageKeysLookup.Lookup(crypto.Keccak256Hash(flop.KicksKey[:]))).To(Equal(flop.KicksMetadata))
			Expect(storageKeysLookup.Lookup(crypto.Keccak256Hash(flop.LiveKey[:]))).To(Equal(flop.LiveMetadata))
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

			Expect(storageRepository.GetFlopBidIdsCalledWith).To(Equal(storageKeysLookup.ContractAddress))
		})

		It("returns error if bid ID lookup fails", func() {
			storageRepository.GetFlopBidIdsError = fakes.FakeError

			_, err := storageKeysLookup.Lookup(fakes.FakeHash)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})
	})

	Describe("bid", func() {
		var fakeBidId string
		var bidBidKey common.Hash

		BeforeEach(func() {
			fakeBidId = "42"
			fakeHexBidId, conversionErr := shared.ConvertIntStringToHex(fakeBidId)

			Expect(conversionErr).NotTo(HaveOccurred())

			bidBidKey = common.BytesToHash(crypto.Keccak256(common.FromHex(fakeHexBidId + flop.BidsIndex)))
			storageRepository.FlopBidIds = []string{fakeBidId}
		})

		It("returns value metadata for bid bid", func() {
			expectedMetadata := utils.StorageValueMetadata{
				Name: mcdStorage.BidBid,
				Keys: map[utils.Key]string{constants.BidId: fakeBidId},
				Type: utils.Uint256,
			}
			Expect(storageKeysLookup.Lookup(bidBidKey)).To(Equal(expectedMetadata))
		})

		It("returns value metadata for bid lot", func() {
			bidLotKey := storage.GetIncrementedKey(bidBidKey, 1)
			expectedMetadata := utils.StorageValueMetadata{
				Name: mcdStorage.BidLot,
				Keys: map[utils.Key]string{constants.BidId: fakeBidId},
				Type: utils.Uint256,
			}
			Expect(storageKeysLookup.Lookup(bidLotKey)).To(Equal(expectedMetadata))
		})

		It("returns value metadata for bid guy + tic + end packed slot", func() {
			bidGuyKey := storage.GetIncrementedKey(bidBidKey, 2)
			expectedMetadata := utils.StorageValueMetadata{
				Name:        mcdStorage.Packed,
				Keys:        map[utils.Key]string{constants.BidId: fakeBidId},
				Type:        utils.PackedSlot,
				PackedTypes: map[int]utils.ValueType{0: utils.Address, 1: utils.Uint48, 2: utils.Uint48},
				PackedNames: map[int]string{0: mcdStorage.BidGuy, 1: mcdStorage.BidTic, 2: mcdStorage.BidEnd},
			}
			Expect(storageKeysLookup.Lookup(bidGuyKey)).To(Equal(expectedMetadata))
		})
	})
})
