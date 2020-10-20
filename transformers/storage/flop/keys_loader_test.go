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
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/flop"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities/wards"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	vdbStorage "github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
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
		storageKeysLoader = flop.NewKeysLoader(storageRepository, test_data.FlopV101Address())
	})

	It("returns value metadata for static keys", func() {
		mappings, err := storageKeysLoader.LoadMappings()

		Expect(err).NotTo(HaveOccurred())
		Expect(mappings[flop.VatKey]).To(Equal(flop.VatMetadata))
		Expect(mappings[flop.GemKey]).To(Equal(flop.GemMetadata))
		Expect(mappings[flop.BegKey]).To(Equal(flop.BegMetadata))
		Expect(mappings[flop.PadKey]).To(Equal(flop.PadMetadata))
		Expect(mappings[flop.TTLAndTauKey]).To(Equal(flop.TTLAndTauMetadata))
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
			expectedMetadata := types.ValueMetadata{
				Name: wards.Wards,
				Keys: map[types.Key]string{constants.User: fakes.FakeAddress.Hex()},
				Type: types.Uint256,
			}

			mappings, err := storageKeysLoader.LoadMappings()
			Expect(err).NotTo(HaveOccurred())

			Expect(storageRepository.GetWardsKeysCalledWith).To(Equal(test_data.FlopV101Address()))
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
				storageRepository.GetFlopBidIDsError = fakes.FakeError

				_, err := storageKeysLoader.LoadMappings()

				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(fakes.FakeError))
			})
		})

		Describe("when getting flop bid IDs succeeds", func() {
			var (
				fakeBidID string
				bidBidKey common.Hash
				mappings  map[common.Hash]types.ValueMetadata
			)

			BeforeEach(func() {
				fakeBidID = "42"
				fakeHexBidID, conversionErr := shared.ConvertIntStringToHex(fakeBidID)

				Expect(conversionErr).NotTo(HaveOccurred())

				bidBidKey = common.BytesToHash(crypto.Keccak256(common.FromHex(fakeHexBidID + flop.BidsIndex)))
				storageRepository.FlopBidIDs = []string{fakeBidID}
				var err error
				mappings, err = storageKeysLoader.LoadMappings()
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns value metadata for bid bid", func() {
				expectedMetadata := types.ValueMetadata{
					Name: mcdStorage.BidBid,
					Keys: map[types.Key]string{constants.BidId: fakeBidID},
					Type: types.Uint256,
				}

				Expect(mappings[bidBidKey]).To(Equal(expectedMetadata))
			})

			It("returns value metadata for bid lot", func() {
				bidLotKey := vdbStorage.GetIncrementedKey(bidBidKey, 1)
				expectedMetadata := types.ValueMetadata{
					Name: mcdStorage.BidLot,
					Keys: map[types.Key]string{constants.BidId: fakeBidID},
					Type: types.Uint256,
				}

				Expect(mappings[bidLotKey]).To(Equal(expectedMetadata))
			})

			It("returns value metadata for bid guy + tic + end packed slot", func() {
				bidGuyKey := vdbStorage.GetIncrementedKey(bidBidKey, 2)
				expectedMetadata := types.ValueMetadata{
					Name:        mcdStorage.Packed,
					Keys:        map[types.Key]string{constants.BidId: fakeBidID},
					Type:        types.PackedSlot,
					PackedTypes: map[int]types.ValueType{0: types.Address, 1: types.Uint48, 2: types.Uint48},
					PackedNames: map[int]string{0: mcdStorage.BidGuy, 1: mcdStorage.BidTic, 2: mcdStorage.BidEnd},
				}

				Expect(mappings[bidGuyKey]).To(Equal(expectedMetadata))
			})
		})
	})
})
