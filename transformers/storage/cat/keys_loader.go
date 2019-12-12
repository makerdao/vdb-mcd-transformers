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
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	vdbStorage "github.com/makerdao/vulcanizedb/libraries/shared/storage"
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
	// wards takes up index 0
	IlksMappingIndex = vdbStorage.IndexOne // bytes32 => flip address; chop (ray), lump (wad) uint256

	LiveKey      = common.HexToHash(vdbStorage.IndexTwo)
	LiveMetadata = vdbStorage.GetStorageValueMetadata(Live, nil, vdbStorage.Uint256)

	VatKey      = common.HexToHash(vdbStorage.IndexThree)
	VatMetadata = vdbStorage.GetStorageValueMetadata(Vat, nil, vdbStorage.Address)

	VowKey      = common.HexToHash(vdbStorage.IndexFour)
	VowMetadata = vdbStorage.GetStorageValueMetadata(Vow, nil, vdbStorage.Address)
)

type keysLoader struct {
	storageRepository mcdStorage.IMakerStorageRepository
}

func NewKeysLoader(storageRepository mcdStorage.IMakerStorageRepository) storage.KeysLoader {
	return &keysLoader{storageRepository: storageRepository}
}

func (loader *keysLoader) SetDB(db *postgres.DB) {
	loader.storageRepository.SetDB(db)
}

func (loader *keysLoader) LoadMappings() (map[common.Hash]vdbStorage.StorageValueMetadata, error) {
	mappings := loadStaticMappings()
	return loader.addIlkKeys(mappings)
}

func (loader *keysLoader) addIlkKeys(mappings map[common.Hash]vdbStorage.StorageValueMetadata) (map[common.Hash]vdbStorage.StorageValueMetadata, error) {
	ilks, err := loader.storageRepository.GetIlks()
	if err != nil {
		return nil, err
	}
	for _, ilk := range ilks {
		mappings[getIlkFlipKey(ilk)] = getIlkFlipMetadata(ilk)
		mappings[getIlkChopKey(ilk)] = getIlkChopMetadata(ilk)
		mappings[getIlkLumpKey(ilk)] = getIlkLumpMetadata(ilk)
	}
	return mappings, nil
}

func loadStaticMappings() map[common.Hash]vdbStorage.StorageValueMetadata {
	mappings := make(map[common.Hash]vdbStorage.StorageValueMetadata)
	mappings[LiveKey] = LiveMetadata
	mappings[VatKey] = VatMetadata
	mappings[VowKey] = VowMetadata
	return mappings
}

func getIlkFlipKey(ilk string) common.Hash {
	return vdbStorage.GetStorageKeyForMapping(IlksMappingIndex, ilk)
}

func getIlkFlipMetadata(ilk string) vdbStorage.StorageValueMetadata {
	keys := map[vdbStorage.Key]string{constants.Ilk: ilk}
	return vdbStorage.GetStorageValueMetadata(IlkFlip, keys, vdbStorage.Address)
}

func getIlkChopKey(ilk string) common.Hash {
	return vdbStorage.GetIncrementedStorageKey(getIlkFlipKey(ilk), 1)
}

func getIlkChopMetadata(ilk string) vdbStorage.StorageValueMetadata {
	keys := map[vdbStorage.Key]string{constants.Ilk: ilk}
	return vdbStorage.GetStorageValueMetadata(IlkChop, keys, vdbStorage.Uint256)
}

func getIlkLumpKey(ilk string) common.Hash {
	return vdbStorage.GetIncrementedStorageKey(getIlkFlipKey(ilk), 2)
}

func getIlkLumpMetadata(ilk string) vdbStorage.StorageValueMetadata {
	keys := map[vdbStorage.Key]string{constants.Ilk: ilk}
	return vdbStorage.GetStorageValueMetadata(IlkLump, keys, vdbStorage.Uint256)
}
