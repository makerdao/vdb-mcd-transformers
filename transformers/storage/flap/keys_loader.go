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

package flap

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	vdbStorage "github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

var (
	BidsIndex = vdbStorage.IndexOne

	VatStorageKey = common.HexToHash(vdbStorage.IndexTwo)
	VatMetadata   = vdbStorage.GetValueMetadata(mcdStorage.Vat, nil, vdbStorage.Address)

	GemStorageKey = common.HexToHash(vdbStorage.IndexThree)
	GemMetadata   = vdbStorage.GetValueMetadata(mcdStorage.Gem, nil, vdbStorage.Address)

	BegStorageKey = common.HexToHash(vdbStorage.IndexFour)
	BegMetadata   = vdbStorage.GetValueMetadata(mcdStorage.Beg, nil, vdbStorage.Uint256)

	TtlAndTauStorageKey = common.HexToHash(vdbStorage.IndexFive)
	ttlAndTauTypes      = map[int]vdbStorage.ValueType{0: vdbStorage.Uint48, 1: vdbStorage.Uint48}
	ttlAndTauNames      = map[int]string{0: mcdStorage.Ttl, 1: mcdStorage.Tau}
	TtlAndTauMetadata   = vdbStorage.GetValueMetadataForPackedSlot(mcdStorage.Packed, nil, vdbStorage.PackedSlot, ttlAndTauNames, ttlAndTauTypes)

	KicksStorageKey = common.HexToHash(vdbStorage.IndexSix)
	KicksMetadata   = vdbStorage.GetValueMetadata(mcdStorage.Kicks, nil, vdbStorage.Uint256)

	LiveStorageKey = common.HexToHash(vdbStorage.IndexSeven)
	LiveMetadata   = vdbStorage.GetValueMetadata(mcdStorage.Live, nil, vdbStorage.Uint256)
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
	mappings := loadStaticKeys()
	return loader.loadBidKeys(mappings)
}

func (loader *keysLoader) loadBidKeys(mappings map[common.Hash]vdbStorage.ValueMetadata) (map[common.Hash]vdbStorage.ValueMetadata, error) {
	bidIds, getBidIdsErr := loader.storageRepository.GetFlapBidIds(loader.contractAddress)
	if getBidIdsErr != nil {
		return nil, getBidIdsErr
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

func loadStaticKeys() map[common.Hash]vdbStorage.ValueMetadata {
	mappings := make(map[common.Hash]vdbStorage.ValueMetadata)
	mappings[VatStorageKey] = VatMetadata
	mappings[GemStorageKey] = GemMetadata
	mappings[BegStorageKey] = BegMetadata
	mappings[TtlAndTauStorageKey] = TtlAndTauMetadata
	mappings[KicksStorageKey] = KicksMetadata
	mappings[LiveStorageKey] = LiveMetadata
	return mappings
}

func getBidBidKey(bidId string) common.Hash {
	return vdbStorage.GetKeyForMapping(BidsIndex, bidId)
}

func getBidBidMetadata(bidId string) vdbStorage.ValueMetadata {
	return vdbStorage.ValueMetadata{
		Name: mcdStorage.BidBid,
		Keys: map[vdbStorage.Key]string{constants.BidId: bidId},
		Type: vdbStorage.Uint256,
	}
}

func getBidLotKey(bidId string) common.Hash {
	return vdbStorage.GetIncrementedKey(getBidBidKey(bidId), 1) //should this be renamed GetMappingKey?
}

func getBidLotMetadata(bidId string) vdbStorage.ValueMetadata {
	return vdbStorage.ValueMetadata{
		Name: mcdStorage.BidLot,
		Keys: map[vdbStorage.Key]string{constants.BidId: bidId},
		Type: vdbStorage.Uint256,
	}
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
