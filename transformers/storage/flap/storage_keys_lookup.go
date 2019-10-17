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
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	mcdStorage "github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

var (
	BidsIndex = storage.IndexOne

	VatStorageKey = common.HexToHash(storage.IndexTwo)
	VatMetadata   = utils.GetStorageValueMetadata(mcdStorage.Vat, nil, utils.Address)

	GemStorageKey = common.HexToHash(storage.IndexThree)
	GemMetadata   = utils.GetStorageValueMetadata(mcdStorage.Gem, nil, utils.Address)

	BegStorageKey = common.HexToHash(storage.IndexFour)
	BegMetadata   = utils.GetStorageValueMetadata(mcdStorage.Beg, nil, utils.Uint256)

	TtlAndTauStorageKey = common.HexToHash(storage.IndexFive)
	ttlAndTauTypes      = map[int]utils.ValueType{0: utils.Uint48, 1: utils.Uint48}
	ttlAndTauNames      = map[int]string{0: mcdStorage.Ttl, 1: mcdStorage.Tau}
	TtlAndTauMetadata   = utils.GetStorageValueMetadataForPackedSlot(mcdStorage.Packed, nil, utils.PackedSlot, ttlAndTauNames, ttlAndTauTypes)

	KicksStorageKey = common.HexToHash(storage.IndexSix)
	KicksMetadata   = utils.GetStorageValueMetadata(mcdStorage.Kicks, nil, utils.Uint256)

	LiveStorageKey = common.HexToHash(storage.IndexSeven)
	LiveMetadata   = utils.GetStorageValueMetadata(mcdStorage.Live, nil, utils.Uint256)
)

type StorageKeysLookup struct {
	StorageRepository mcdStorage.IMakerStorageRepository
	mappings          map[common.Hash]utils.StorageValueMetadata
	ContractAddress   string
}

func (lookup *StorageKeysLookup) Lookup(key common.Hash) (utils.StorageValueMetadata, error) {
	metadata, ok := lookup.mappings[key]
	if !ok {
		loadErr := lookup.loadMapping()
		if loadErr != nil {
			return utils.StorageValueMetadata{}, loadErr
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

func (lookup *StorageKeysLookup) loadMapping() error {
	lookup.loadStaticKeys()
	err := lookup.loadBidKeys()
	if err != nil {
		return err
	}
	lookup.mappings = storage.AddHashedKeys(lookup.mappings)
	return nil
}

func (lookup *StorageKeysLookup) loadBidKeys() error {
	bidIds, getBidIdsErr := lookup.StorageRepository.GetFlapBidIds(lookup.ContractAddress)
	for _, bidId := range bidIds {
		hexBidId, convertErr := shared.ConvertIntStringToHex(bidId)
		if convertErr != nil {
			return convertErr
		}

		lookup.mappings[getBidBidKey(hexBidId)] = getBidBidMetadata(bidId)
		lookup.mappings[getBidLotKey(hexBidId)] = getBidLotMetadata(bidId)
		lookup.mappings[getBidGuyTicEndKey(hexBidId)] = getBidGuyTicEndMetadata(bidId)
	}

	return getBidIdsErr
}

func getBidBidKey(bidId string) common.Hash {
	return storage.GetMapping(BidsIndex, bidId)
}

func getBidBidMetadata(bidId string) utils.StorageValueMetadata {
	return utils.StorageValueMetadata{
		Name: mcdStorage.BidBid,
		Keys: map[utils.Key]string{constants.BidId: bidId},
		Type: utils.Uint256,
	}
}

func getBidLotKey(bidId string) common.Hash {
	return storage.GetIncrementedKey(getBidBidKey(bidId), 1) //should this be renamed GetMappingKey?
}

func getBidLotMetadata(bidId string) utils.StorageValueMetadata {
	return utils.StorageValueMetadata{
		Name: mcdStorage.BidLot,
		Keys: map[utils.Key]string{constants.BidId: bidId},
		Type: utils.Uint256,
	}
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

func (lookup *StorageKeysLookup) loadStaticKeys() {
	mappings := make(map[common.Hash]utils.StorageValueMetadata)
	mappings[VatStorageKey] = VatMetadata
	mappings[GemStorageKey] = GemMetadata
	mappings[BegStorageKey] = BegMetadata
	mappings[TtlAndTauStorageKey] = TtlAndTauMetadata
	mappings[KicksStorageKey] = KicksMetadata
	mappings[LiveStorageKey] = LiveMetadata
	lookup.mappings = mappings
}
