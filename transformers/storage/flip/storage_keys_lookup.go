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
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	mcdStorage "github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

var (
	BidsMappingIndex = storage.IndexOne

	VatKey      = common.HexToHash(storage.IndexTwo)
	VatMetadata = utils.GetStorageValueMetadata(mcdStorage.Vat, nil, utils.Address)

	IlkKey      = common.HexToHash(storage.IndexThree)
	IlkMetadata = utils.GetStorageValueMetadata(mcdStorage.Ilk, nil, utils.Bytes32)

	BegKey      = common.HexToHash(storage.IndexFour)
	BegMetadata = utils.GetStorageValueMetadata(mcdStorage.Beg, nil, utils.Uint256)

	TtlAndTauStorageKey = common.HexToHash(storage.IndexFive)
	ttlAndTauTypes      = map[int]utils.ValueType{0: utils.Uint48, 1: utils.Uint48}
	ttlAndTauNames      = map[int]string{0: mcdStorage.Ttl, 1: mcdStorage.Tau}
	TtlAndTauMetadata   = utils.GetStorageValueMetadataForPackedSlot(mcdStorage.Packed, nil, utils.PackedSlot, ttlAndTauNames, ttlAndTauTypes)

	KicksKey      = common.HexToHash(storage.IndexSix)
	KicksMetadata = utils.GetStorageValueMetadata(mcdStorage.Kicks, nil, utils.Uint256)
)

type StorageKeysLookup struct {
	StorageRepository mcdStorage.IMakerStorageRepository
	ContractAddress   string
	mappings          map[common.Hash]utils.StorageValueMetadata
}

func (lookup StorageKeysLookup) Lookup(key common.Hash) (utils.StorageValueMetadata, error) {
	metadata, ok := lookup.mappings[key]
	if !ok {
		err := lookup.loadMappings()
		if err != nil {
			return metadata, err
		}
		metadata, ok = lookup.mappings[key]
		if !ok {
			return metadata, utils.ErrStorageKeyNotFound{Key: key.Hex()}
		}
	}
	return metadata, nil
}

func (lookup *StorageKeysLookup) SetDB(db *postgres.DB) {
	lookup.StorageRepository.SetDB(db)
}

func (lookup *StorageKeysLookup) loadMappings() error {
	lookup.mappings = loadStaticMappings()
	err := lookup.loadBidKeys()
	if err != nil {
		return err
	}
	lookup.mappings = storage.AddHashedKeys(lookup.mappings)
	return nil
}

func loadStaticMappings() map[common.Hash]utils.StorageValueMetadata {
	mappings := make(map[common.Hash]utils.StorageValueMetadata)
	mappings[VatKey] = VatMetadata
	mappings[IlkKey] = IlkMetadata
	mappings[BegKey] = BegMetadata
	mappings[TtlAndTauStorageKey] = TtlAndTauMetadata
	mappings[KicksKey] = KicksMetadata
	return mappings
}

func (lookup *StorageKeysLookup) loadBidKeys() error {
	bidIds, err := lookup.StorageRepository.GetFlipBidIds(lookup.ContractAddress)
	if err != nil {
		return err
	}
	for _, bidId := range bidIds {
		hexBidId, err := shared.ConvertIntStringToHex(bidId)
		if err != nil {
			return err
		}
		lookup.mappings[getBidBidKey(hexBidId)] = getBidBidMetadata(bidId)
		lookup.mappings[getBidLotKey(hexBidId)] = getBidLotMetadata(bidId)
		lookup.mappings[getBidGuyTicEndKey(hexBidId)] = getBidGuyTicEndMetadata(bidId)
		lookup.mappings[getBidUsrKey(hexBidId)] = getBidUsrMetadata(bidId)
		lookup.mappings[getBidGalKey(hexBidId)] = getBidGalMetadata(bidId)
		lookup.mappings[getBidTabKey(hexBidId)] = getBidTabMetadata(bidId)
	}
	return nil
}

func getBidBidKey(hexBidId string) common.Hash {
	return storage.GetMapping(BidsMappingIndex, hexBidId)
}

func getBidBidMetadata(bidId string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.BidId: bidId}
	return utils.GetStorageValueMetadata(mcdStorage.BidBid, keys, utils.Uint256)
}

func getBidLotKey(hexBidId string) common.Hash {
	return storage.GetIncrementedKey(getBidBidKey(hexBidId), 1)
}

func getBidLotMetadata(bidId string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.BidId: bidId}
	return utils.GetStorageValueMetadata(mcdStorage.BidLot, keys, utils.Uint256)
}

func getBidGuyTicEndKey(hexBidId string) common.Hash {
	return storage.GetIncrementedKey(getBidBidKey(hexBidId), 2)
}

func getBidGuyTicEndMetadata(bidId string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.BidId: bidId}
	packedTypes := map[int]utils.ValueType{0: utils.Address, 1: utils.Uint48, 2: utils.Uint48}
	packedNames := map[int]string{0: mcdStorage.BidGuy, 1: mcdStorage.BidTic, 2: mcdStorage.BidEnd}
	return utils.GetStorageValueMetadataForPackedSlot(mcdStorage.Packed, keys, utils.PackedSlot, packedNames, packedTypes)
}

func getBidUsrKey(hexBidId string) common.Hash {
	return storage.GetIncrementedKey(getBidBidKey(hexBidId), 3)
}

func getBidUsrMetadata(bidId string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.BidId: bidId}
	return utils.GetStorageValueMetadata(mcdStorage.BidUsr, keys, utils.Address)
}

func getBidGalKey(hexBidId string) common.Hash {
	return storage.GetIncrementedKey(getBidBidKey(hexBidId), 4)
}

func getBidGalMetadata(bidId string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.BidId: bidId}
	return utils.GetStorageValueMetadata(mcdStorage.BidGal, keys, utils.Address)
}

func getBidTabKey(hexBidId string) common.Hash {
	return storage.GetIncrementedKey(getBidBidKey(hexBidId), 5)
}

func getBidTabMetadata(bidId string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.BidId: bidId}
	return utils.GetStorageValueMetadata(mcdStorage.BidTab, keys, utils.Uint256)
}
