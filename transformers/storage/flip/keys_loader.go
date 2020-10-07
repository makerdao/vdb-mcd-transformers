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
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities/wards"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	vdbStorage "github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

var (
	BidsMappingIndex = vdbStorage.IndexOne

	VatKey      = common.HexToHash(vdbStorage.IndexTwo)
	VatMetadata = types.GetValueMetadata(mcdStorage.Vat, nil, types.Address)

	IlkKey      = common.HexToHash(vdbStorage.IndexThree)
	IlkMetadata = types.GetValueMetadata(mcdStorage.Ilk, nil, types.Bytes32)

	BegKey      = common.HexToHash(vdbStorage.IndexFour)
	BegMetadata = types.GetValueMetadata(mcdStorage.Beg, nil, types.Uint256)

	TTLAndTauStorageKey = common.HexToHash(vdbStorage.IndexFive)
	ttlAndTauTypes      = map[int]types.ValueType{0: types.Uint48, 1: types.Uint48}
	ttlAndTauNames      = map[int]string{0: mcdStorage.Ttl, 1: mcdStorage.Tau}
	TTLAndTauMetadata   = types.GetValueMetadataForPackedSlot(mcdStorage.Packed, nil, types.PackedSlot, ttlAndTauNames, ttlAndTauTypes)

	KicksKey      = common.HexToHash(vdbStorage.IndexSix)
	KicksMetadata = types.GetValueMetadata(mcdStorage.Kicks, nil, types.Uint256)

	CatKey      = common.HexToHash(vdbStorage.IndexSeven)
	CatMetadata = types.GetValueMetadata(mcdStorage.Cat, nil, types.Address)
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

func (loader *keysLoader) LoadMappings() (map[common.Hash]types.ValueMetadata, error) {
	mappings := loadStaticMappings()
	mappings, wardsErr := loader.loadWardsKeys(mappings)
	if wardsErr != nil {
		return nil, fmt.Errorf("error adding ward keys to flip keys loader: %w", wardsErr)
	}
	mappings, bidErr := loader.loadBidKeys(mappings)
	if bidErr != nil {
		return nil, fmt.Errorf("error adding bid keys to flip keys loader: %w", bidErr)
	}
	return mappings, nil
}

func (loader *keysLoader) loadWardsKeys(mappings map[common.Hash]types.ValueMetadata) (map[common.Hash]types.ValueMetadata, error) {
	addresses, err := loader.storageRepository.GetWardsAddresses(loader.contractAddress)
	if err != nil {
		return nil, fmt.Errorf("error getting wards addresses: %w", err)
	}
	return wards.AddWardsKeys(mappings, addresses)
}

func (loader *keysLoader) loadBidKeys(mappings map[common.Hash]types.ValueMetadata) (map[common.Hash]types.ValueMetadata, error) {
	bidIDs, bidErr := loader.storageRepository.GetFlipBidIDs(loader.contractAddress)
	if bidErr != nil {
		return nil, fmt.Errorf("error getting flip bid IDs: %w", bidErr)
	}
	for _, bidID := range bidIDs {
		hexBidID, convertErr := shared.ConvertIntStringToHex(bidID)
		if convertErr != nil {
			return nil, fmt.Errorf("error converting int string to hex: %w", convertErr)
		}
		mappings[getBidBidKey(hexBidID)] = getBidBidMetadata(bidID)
		mappings[getBidLotKey(hexBidID)] = getBidLotMetadata(bidID)
		mappings[getBidGuyTicEndKey(hexBidID)] = getBidGuyTicEndMetadata(bidID)
		mappings[getBidUsrKey(hexBidID)] = getBidUsrMetadata(bidID)
		mappings[getBidGalKey(hexBidID)] = getBidGalMetadata(bidID)
		mappings[getBidTabKey(hexBidID)] = getBidTabMetadata(bidID)
	}
	return mappings, nil
}

func loadStaticMappings() map[common.Hash]types.ValueMetadata {
	mappings := make(map[common.Hash]types.ValueMetadata)
	mappings[VatKey] = VatMetadata
	mappings[IlkKey] = IlkMetadata
	mappings[BegKey] = BegMetadata
	mappings[TTLAndTauStorageKey] = TTLAndTauMetadata
	mappings[KicksKey] = KicksMetadata
	mappings[CatKey] = CatMetadata
	return mappings
}

func getBidBidKey(hexBidID string) common.Hash {
	return vdbStorage.GetKeyForMapping(BidsMappingIndex, hexBidID)
}

func getBidBidMetadata(bidID string) types.ValueMetadata {
	keys := map[types.Key]string{constants.BidId: bidID}
	return types.GetValueMetadata(mcdStorage.BidBid, keys, types.Uint256)
}

func getBidLotKey(hexBidID string) common.Hash {
	return vdbStorage.GetIncrementedKey(getBidBidKey(hexBidID), 1)
}

func getBidLotMetadata(bidID string) types.ValueMetadata {
	keys := map[types.Key]string{constants.BidId: bidID}
	return types.GetValueMetadata(mcdStorage.BidLot, keys, types.Uint256)
}

func getBidGuyTicEndKey(hexBidID string) common.Hash {
	return vdbStorage.GetIncrementedKey(getBidBidKey(hexBidID), 2)
}

func getBidGuyTicEndMetadata(bidID string) types.ValueMetadata {
	keys := map[types.Key]string{constants.BidId: bidID}
	packedTypes := map[int]types.ValueType{0: types.Address, 1: types.Uint48, 2: types.Uint48}
	packedNames := map[int]string{0: mcdStorage.BidGuy, 1: mcdStorage.BidTic, 2: mcdStorage.BidEnd}
	return types.GetValueMetadataForPackedSlot(mcdStorage.Packed, keys, types.PackedSlot, packedNames, packedTypes)
}

func getBidUsrKey(hexBidID string) common.Hash {
	return vdbStorage.GetIncrementedKey(getBidBidKey(hexBidID), 3)
}

func getBidUsrMetadata(bidID string) types.ValueMetadata {
	keys := map[types.Key]string{constants.BidId: bidID}
	return types.GetValueMetadata(mcdStorage.BidUsr, keys, types.Address)
}

func getBidGalKey(hexBidID string) common.Hash {
	return vdbStorage.GetIncrementedKey(getBidBidKey(hexBidID), 4)
}

func getBidGalMetadata(bidID string) types.ValueMetadata {
	keys := map[types.Key]string{constants.BidId: bidID}
	return types.GetValueMetadata(mcdStorage.BidGal, keys, types.Address)
}

func getBidTabKey(hexBidID string) common.Hash {
	return vdbStorage.GetIncrementedKey(getBidBidKey(hexBidID), 5)
}

func getBidTabMetadata(bidID string) types.ValueMetadata {
	keys := map[types.Key]string{constants.BidId: bidID}
	return types.GetValueMetadata(mcdStorage.BidTab, keys, types.Uint256)
}
