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

package cat

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities/wards"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	vdbStorage "github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

const (
	Live = "live"
	Vat  = "vat"
	Vow  = "vow"

	IlkFlip = "flip"
	IlkChop = "chop"
	IlkLump = "lump"
)

var (
	IlksMappingIndex = vdbStorage.IndexOne // bytes32 => flip address; chop (ray), lump (wad) uint256

	LiveKey      = common.HexToHash(vdbStorage.IndexTwo)
	LiveMetadata = types.GetValueMetadata(Live, nil, types.Uint256)

	VatKey      = common.HexToHash(vdbStorage.IndexThree)
	VatMetadata = types.GetValueMetadata(Vat, nil, types.Address)

	VowKey      = common.HexToHash(vdbStorage.IndexFour)
	VowMetadata = types.GetValueMetadata(Vow, nil, types.Address)
)

type keysLoader struct {
	storageRepository mcdStorage.IMakerStorageRepository
	contractAddress   string
}

func NewKeysLoader(storageRepository mcdStorage.IMakerStorageRepository, contractAddress string) storage.KeysLoader {
	return &keysLoader{storageRepository: storageRepository, contractAddress: contractAddress}
}

func (loader *keysLoader) SetDB(db *postgres.DB) {
	loader.storageRepository.SetDB(db)
}

func (loader *keysLoader) LoadMappings() (map[common.Hash]types.ValueMetadata, error) {
	mappings := loadStaticMappings()
	mappings, ilkErr := loader.addIlkKeys(mappings)
	if ilkErr != nil {
		return nil, fmt.Errorf("error adding ilk keys to cat keys loader: %w", ilkErr)
	}
	mappings, wardsErr := loader.addWardsKeys(mappings)
	if wardsErr != nil {
		return nil, fmt.Errorf("error adding wards keys to cat keys loader: %w", wardsErr)
	}
	return mappings, nil
}

func (loader *keysLoader) addIlkKeys(mappings map[common.Hash]types.ValueMetadata) (map[common.Hash]types.ValueMetadata, error) {
	ilks, err := loader.storageRepository.GetIlks()
	if err != nil {
		return nil, fmt.Errorf("error getting ilks: %w", err)
	}
	for _, ilk := range ilks {
		mappings[getIlkFlipKey(ilk)] = getIlkFlipMetadata(ilk)
		mappings[getIlkChopKey(ilk)] = getIlkChopMetadata(ilk)
		mappings[getIlkLumpKey(ilk)] = getIlkLumpMetadata(ilk)
	}
	return mappings, nil
}

func (loader *keysLoader) addWardsKeys(mappings map[common.Hash]types.ValueMetadata) (map[common.Hash]types.ValueMetadata, error) {
	addresses, err := loader.storageRepository.GetWardsAddresses(loader.contractAddress)
	if err != nil {
		return nil, fmt.Errorf("error getting wards addresses: %w", err)
	}
	return wards.AddWardsKeys(mappings, addresses)
}

func loadStaticMappings() map[common.Hash]types.ValueMetadata {
	mappings := make(map[common.Hash]types.ValueMetadata)
	mappings[LiveKey] = LiveMetadata
	mappings[VatKey] = VatMetadata
	mappings[VowKey] = VowMetadata
	return mappings
}

func getIlkFlipKey(ilk string) common.Hash {
	return vdbStorage.GetKeyForMapping(IlksMappingIndex, ilk)
}

func getIlkFlipMetadata(ilk string) types.ValueMetadata {
	keys := map[types.Key]string{constants.Ilk: ilk}
	return types.GetValueMetadata(IlkFlip, keys, types.Address)
}

func getIlkChopKey(ilk string) common.Hash {
	return vdbStorage.GetIncrementedKey(getIlkFlipKey(ilk), 1)
}

func getIlkChopMetadata(ilk string) types.ValueMetadata {
	keys := map[types.Key]string{constants.Ilk: ilk}
	return types.GetValueMetadata(IlkChop, keys, types.Uint256)
}

func getIlkLumpKey(ilk string) common.Hash {
	return vdbStorage.GetIncrementedKey(getIlkFlipKey(ilk), 2)
}

func getIlkLumpMetadata(ilk string) types.ValueMetadata {
	keys := map[types.Key]string{constants.Ilk: ilk}
	return types.GetValueMetadata(IlkLump, keys, types.Uint256)
}
