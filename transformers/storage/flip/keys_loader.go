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

package flip

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
	BidsMappingIndex = vdbStorage.IndexOne

	VatKey      = common.HexToHash(vdbStorage.IndexTwo)
	VatMetadata = vdbStorage.GetStorageValueMetadata(mcdStorage.Vat, nil, vdbStorage.Address)

	IlkKey      = common.HexToHash(vdbStorage.IndexThree)
	IlkMetadata = vdbStorage.GetStorageValueMetadata(mcdStorage.Ilk, nil, vdbStorage.Bytes32)

	BegKey      = common.HexToHash(vdbStorage.IndexFour)
	BegMetadata = vdbStorage.GetStorageValueMetadata(mcdStorage.Beg, nil, vdbStorage.Uint256)

	TtlAndTauStorageKey = common.HexToHash(vdbStorage.IndexFive)
	ttlAndTauTypes      = map[int]vdbStorage.ValueType{0: vdbStorage.Uint48, 1: vdbStorage.Uint48}
	ttlAndTauNames      = map[int]string{0: mcdStorage.Ttl, 1: mcdStorage.Tau}
	TtlAndTauMetadata   = vdbStorage.GetStorageValueMetadataForPackedSlot(mcdStorage.Packed, nil, vdbStorage.PackedSlot, ttlAndTauNames, ttlAndTauTypes)

	KicksKey      = common.HexToHash(vdbStorage.IndexSix)
	KicksMetadata = vdbStorage.GetStorageValueMetadata(mcdStorage.Kicks, nil, vdbStorage.Uint256)
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

func (loader *keysLoader) LoadMappings() (map[common.Hash]vdbStorage.StorageValueMetadata, error) {
	mappings := loadStaticMappings()
	return loader.loadBidKeys(mappings)
}

func (loader *keysLoader) loadBidKeys(mappings map[common.Hash]vdbStorage.StorageValueMetadata) (map[common.Hash]vdbStorage.StorageValueMetadata, error) {
	bidIds, bidErr := loader.storageRepository.GetFlipBidIds(loader.contractAddress)
	if bidErr != nil {
		return nil, bidErr
	}
	for _, bidId := range bidIds {
		hexBidId, convertErr := shared.ConvertIntStringToHex(bidId)
		if convertErr != nil {
			return nil, convertErr
		}
		mappings[getBidBidKey(hexBidId)] = getBidBidMetadata(bidId)
		mappings[getBidLotKey(hexBidId)] = getBidLotMetadata(bidId)
		mappings[getBidGuyTicEndKey(hexBidId)] = getBidGuyTicEndMetadata(bidId)
		mappings[getBidUsrKey(hexBidId)] = getBidUsrMetadata(bidId)
		mappings[getBidGalKey(hexBidId)] = getBidGalMetadata(bidId)
		mappings[getBidTabKey(hexBidId)] = getBidTabMetadata(bidId)
	}
	return mappings, nil
}

func loadStaticMappings() map[common.Hash]vdbStorage.StorageValueMetadata {
	mappings := make(map[common.Hash]vdbStorage.StorageValueMetadata)
	mappings[VatKey] = VatMetadata
	mappings[IlkKey] = IlkMetadata
	mappings[BegKey] = BegMetadata
	mappings[TtlAndTauStorageKey] = TtlAndTauMetadata
	mappings[KicksKey] = KicksMetadata
	return mappings
}

func getBidBidKey(hexBidId string) common.Hash {
	return vdbStorage.GetStorageKeyForMapping(BidsMappingIndex, hexBidId)
}

func getBidBidMetadata(bidId string) vdbStorage.StorageValueMetadata {
	keys := map[vdbStorage.Key]string{constants.BidId: bidId}
	return vdbStorage.GetStorageValueMetadata(mcdStorage.BidBid, keys, vdbStorage.Uint256)
}

func getBidLotKey(hexBidId string) common.Hash {
	return vdbStorage.GetIncrementedStorageKey(getBidBidKey(hexBidId), 1)
}

func getBidLotMetadata(bidId string) vdbStorage.StorageValueMetadata {
	keys := map[vdbStorage.Key]string{constants.BidId: bidId}
	return vdbStorage.GetStorageValueMetadata(mcdStorage.BidLot, keys, vdbStorage.Uint256)
}

func getBidGuyTicEndKey(hexBidId string) common.Hash {
	return vdbStorage.GetIncrementedStorageKey(getBidBidKey(hexBidId), 2)
}

func getBidGuyTicEndMetadata(bidId string) vdbStorage.StorageValueMetadata {
	keys := map[vdbStorage.Key]string{constants.BidId: bidId}
	packedTypes := map[int]vdbStorage.ValueType{0: vdbStorage.Address, 1: vdbStorage.Uint48, 2: vdbStorage.Uint48}
	packedNames := map[int]string{0: mcdStorage.BidGuy, 1: mcdStorage.BidTic, 2: mcdStorage.BidEnd}
	return vdbStorage.GetStorageValueMetadataForPackedSlot(mcdStorage.Packed, keys, vdbStorage.PackedSlot, packedNames, packedTypes)
}

func getBidUsrKey(hexBidId string) common.Hash {
	return vdbStorage.GetIncrementedStorageKey(getBidBidKey(hexBidId), 3)
}

func getBidUsrMetadata(bidId string) vdbStorage.StorageValueMetadata {
	keys := map[vdbStorage.Key]string{constants.BidId: bidId}
	return vdbStorage.GetStorageValueMetadata(mcdStorage.BidUsr, keys, vdbStorage.Address)
}

func getBidGalKey(hexBidId string) common.Hash {
	return vdbStorage.GetIncrementedStorageKey(getBidBidKey(hexBidId), 4)
}

func getBidGalMetadata(bidId string) vdbStorage.StorageValueMetadata {
	keys := map[vdbStorage.Key]string{constants.BidId: bidId}
	return vdbStorage.GetStorageValueMetadata(mcdStorage.BidGal, keys, vdbStorage.Address)
}

func getBidTabKey(hexBidId string) common.Hash {
	return vdbStorage.GetIncrementedStorageKey(getBidBidKey(hexBidId), 5)
}

func getBidTabMetadata(bidId string) vdbStorage.StorageValueMetadata {
	keys := map[vdbStorage.Key]string{constants.BidId: bidId}
	return vdbStorage.GetStorageValueMetadata(mcdStorage.BidTab, keys, vdbStorage.Uint256)
}
