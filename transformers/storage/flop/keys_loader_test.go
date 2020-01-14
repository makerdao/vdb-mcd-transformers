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
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/flop"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities/wards"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	vdbStorage "github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Flop storage keys loader", func() {
	var (
		storageRepository *test_helpers.MockMakerStorageRepository
		storageKeysLoader storage.KeysLoader
	)

	BeforeEach(func() {
		storageRepository = &test_helpers.MockMakerStorageRepository{}
		storageKeysLoader = flop.NewKeysLoader(storageRepository, test_data.FlopAddress())
	})

	It("returns value metadata for static keys", func() {
		mappings, err := storageKeysLoader.LoadMappings()

		Expect(err).NotTo(HaveOccurred())
		Expect(mappings[flop.VatKey]).To(Equal(flop.VatMetadata))
		Expect(mappings[flop.GemKey]).To(Equal(flop.GemMetadata))
		Expect(mappings[flop.BegKey]).To(Equal(flop.BegMetadata))
		Expect(mappings[flop.PadKey]).To(Equal(flop.PadMetadata))
		Expect(mappings[flop.TtlAndTauKey]).To(Equal(flop.TtlAndTauMetadata))
		Expect(mappings[flop.KicksKey]).To(Equal(flop.KicksMetadata))
		Expect(mappings[flop.LiveKey]).To(Equal(flop.LiveMetadata))
		Expect(mappings[flop.VowKey]).To(Equal(flop.VowMetadata))
	})

	Describe("wards", func() {
		It("returns value metadata for wards", func() {
			wardsUser := fakes.FakeAddress.Hex()
			storageRepository.WardsKeys = []string{wardsUser}
			paddedWardsUser := "0x000000000000000000000000" + wardsUser[2:]
			wardsKey := common.BytesToHash(crypto.Keccak256(common.FromHex(paddedWardsUser + wards.WardsMappingIndex)))
			expectedMetadata := vdbStorage.ValueMetadata{
				Name: wards.Wards,
				Keys: map[vdbStorage.Key]string{constants.User: fakes.FakeAddress.Hex()},
				Type: vdbStorage.Uint256,
			}

			mappings, err := storageKeysLoader.LoadMappings()
			Expect(err).NotTo(HaveOccurred())

			Expect(storageRepository.GetWardsKeysCalledWith).To(Equal(test_data.FlopAddress()))
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

	Describe("bid", func() {
		Describe("when getting flop bid IDs fails", func() {
			It("returns error", func() {
				storageRepository.GetFlopBidIdsError = fakes.FakeError

				_, err := storageKeysLoader.LoadMappings()

				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(fakes.FakeError))
			})
		})

		Describe("when getting flop bid IDs succeeds", func() {
			var (
				fakeBidId string
				bidBidKey common.Hash
				mappings  map[common.Hash]vdbStorage.ValueMetadata
			)

			BeforeEach(func() {
				fakeBidId = "42"
				fakeHexBidId, conversionErr := shared.ConvertIntStringToHex(fakeBidId)

				Expect(conversionErr).NotTo(HaveOccurred())

				bidBidKey = common.BytesToHash(crypto.Keccak256(common.FromHex(fakeHexBidId + flop.BidsIndex)))
				storageRepository.FlopBidIds = []string{fakeBidId}
				var err error
				mappings, err = storageKeysLoader.LoadMappings()
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns value metadata for bid bid", func() {
				expectedMetadata := vdbStorage.ValueMetadata{
					Name: mcdStorage.BidBid,
					Keys: map[vdbStorage.Key]string{constants.BidId: fakeBidId},
					Type: vdbStorage.Uint256,
				}

				Expect(mappings[bidBidKey]).To(Equal(expectedMetadata))
			})

			It("returns value metadata for bid lot", func() {
				bidLotKey := vdbStorage.GetIncrementedKey(bidBidKey, 1)
				expectedMetadata := vdbStorage.ValueMetadata{
					Name: mcdStorage.BidLot,
					Keys: map[vdbStorage.Key]string{constants.BidId: fakeBidId},
					Type: vdbStorage.Uint256,
				}

				Expect(mappings[bidLotKey]).To(Equal(expectedMetadata))
			})

			It("returns value metadata for bid guy + tic + end packed slot", func() {
				bidGuyKey := vdbStorage.GetIncrementedKey(bidBidKey, 2)
				expectedMetadata := vdbStorage.ValueMetadata{
					Name:        mcdStorage.Packed,
					Keys:        map[vdbStorage.Key]string{constants.BidId: fakeBidId},
					Type:        vdbStorage.PackedSlot,
					PackedTypes: map[int]vdbStorage.ValueType{0: vdbStorage.Address, 1: vdbStorage.Uint48, 2: vdbStorage.Uint48},
					PackedNames: map[int]string{0: mcdStorage.BidGuy, 1: mcdStorage.BidTic, 2: mcdStorage.BidEnd},
				}

				Expect(mappings[bidGuyKey]).To(Equal(expectedMetadata))
			})
		})
	})
})
