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

package flop

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities/wards"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	vdbStorage "github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

var (
	BidsIndex = vdbStorage.IndexOne

	VatKey      = common.HexToHash(vdbStorage.IndexTwo)
	VatMetadata = vdbStorage.GetValueMetadata(mcdStorage.Vat, nil, vdbStorage.Address)

	GemKey      = common.HexToHash(vdbStorage.IndexThree)
	GemMetadata = vdbStorage.GetValueMetadata(mcdStorage.Gem, nil, vdbStorage.Address)

	BegKey      = common.HexToHash(vdbStorage.IndexFour)
	BegMetadata = vdbStorage.GetValueMetadata(mcdStorage.Beg, nil, vdbStorage.Uint256)

	PadKey      = common.HexToHash(vdbStorage.IndexFive)
	PadMetadata = vdbStorage.GetValueMetadata(mcdStorage.Pad, nil, vdbStorage.Uint256)

	TtlAndTauKey      = common.HexToHash(vdbStorage.IndexSix)
	ttlAndTauTypes    = map[int]vdbStorage.ValueType{0: vdbStorage.Uint48, 1: vdbStorage.Uint48}
	ttlAndTauNames    = map[int]string{0: mcdStorage.Ttl, 1: mcdStorage.Tau}
	TtlAndTauMetadata = vdbStorage.GetValueMetadataForPackedSlot(mcdStorage.Packed, nil, vdbStorage.PackedSlot, ttlAndTauNames, ttlAndTauTypes)

	KicksKey      = common.HexToHash(vdbStorage.IndexSeven)
	KicksMetadata = vdbStorage.GetValueMetadata(mcdStorage.Kicks, nil, vdbStorage.Uint256)

	LiveKey      = common.HexToHash(vdbStorage.IndexEight)
	LiveMetadata = vdbStorage.GetValueMetadata(mcdStorage.Live, nil, vdbStorage.Uint256)

	VowKey      = common.HexToHash(vdbStorage.IndexNine)
	VowMetadata = vdbStorage.GetValueMetadata(mcdStorage.Vow, nil, vdbStorage.Address)
)

type keysLoader struct {
	storageRepository mcdStorage.IMakerStorageRepository
	contractAddress   string
}

func NewKeysLoader(storageRepository mcdStorage.IMakerStorageRepository, contractAddress string) storage.KeysLoader {
	return &keysLoader{
		storageRepository: storageRepository,
		contractAddress:   contractAddress,
	}
}

func (loader *keysLoader) SetDB(db *postgres.DB) {
	loader.storageRepository.SetDB(db)
}

func (loader *keysLoader) LoadMappings() (map[common.Hash]vdbStorage.ValueMetadata, error) {
	mappings := loadStaticMappings()
	mappings, wardsErr := loader.loadWardsKeys(mappings)
	if wardsErr != nil {
		return nil, wardsErr
	}
	return loader.loadBidKeys(mappings)
}

func (loader *keysLoader) loadWardsKeys(mappings map[common.Hash]vdbStorage.ValueMetadata) (map[common.Hash]vdbStorage.ValueMetadata, error) {
	addresses, err := loader.storageRepository.GetWardsAddresses(loader.contractAddress)
	if err != nil {
		return nil, err
	}
	return wards.AddWardsKeys(mappings, addresses)
}

func (loader *keysLoader) loadBidKeys(mappings map[common.Hash]vdbStorage.ValueMetadata) (map[common.Hash]vdbStorage.ValueMetadata, error) {
	bidIds, getBidsErr := loader.storageRepository.GetFlopBidIds(loader.contractAddress)
	if getBidsErr != nil {
		return nil, getBidsErr
	}
	for _, bidId := range bidIds {
		hexBidId, convertErr := shared.ConvertIntStringToHex(bidId)
		if convertErr != nil {
			return nil, convertErr
		}
		mappings[getBidBidKey(hexBidId)] = getBidBidMetadata(bidId)
		mappings[getBidLotKey(hexBidId)] = getBidLotMetadata(bidId)
		mappings[getBidGuyTicEndKey(hexBidId)] = getBidGuyTicEndMetadata(bidId)
	}
	return mappings, nil
}

func loadStaticMappings() map[common.Hash]vdbStorage.ValueMetadata {
	mappings := make(map[common.Hash]vdbStorage.ValueMetadata)
	mappings[VatKey] = VatMetadata
	mappings[GemKey] = GemMetadata
	mappings[BegKey] = BegMetadata
	mappings[PadKey] = PadMetadata
	mappings[TtlAndTauKey] = TtlAndTauMetadata
	mappings[KicksKey] = KicksMetadata
	mappings[LiveKey] = LiveMetadata
	mappings[VowKey] = VowMetadata
	return mappings
}

func getBidBidKey(hexBidId string) common.Hash {
	return vdbStorage.GetKeyForMapping(BidsIndex, hexBidId)
}

func getBidBidMetadata(bidId string) vdbStorage.ValueMetadata {
	keys := map[vdbStorage.Key]string{constants.BidId: bidId}
	return vdbStorage.GetValueMetadata(mcdStorage.BidBid, keys, vdbStorage.Uint256)
}

func getBidLotKey(hexBidId string) common.Hash {
	return vdbStorage.GetIncrementedKey(getBidBidKey(hexBidId), 1)
}

func getBidLotMetadata(bidId string) vdbStorage.ValueMetadata {
	keys := map[vdbStorage.Key]string{constants.BidId: bidId}
	return vdbStorage.GetValueMetadata(mcdStorage.BidLot, keys, vdbStorage.Uint256)
}

func getBidGuyTicEndKey(hexBidId string) common.Hash {
	return vdbStorage.GetIncrementedKey(getBidBidKey(hexBidId), 2)
}

func getBidGuyTicEndMetadata(bidId string) vdbStorage.ValueMetadata {
	keys := map[vdbStorage.Key]string{constants.BidId: bidId}
	packedTypes := map[int]vdbStorage.ValueType{0: vdbStorage.Address, 1: vdbStorage.Uint48, 2: vdbStorage.Uint48}
	packedNames := map[int]string{0: mcdStorage.BidGuy, 1: mcdStorage.BidTic, 2: mcdStorage.BidEnd}
	return vdbStorage.GetValueMetadataForPackedSlot(mcdStorage.Packed, keys, vdbStorage.PackedSlot, packedNames, packedTypes)
}
