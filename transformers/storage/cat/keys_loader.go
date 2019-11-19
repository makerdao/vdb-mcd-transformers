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
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/utils"
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
	IlksMappingIndex = utils.IndexOne // bytes32 => flip address; chop (ray), lump (wad) uint256

	LiveKey      = common.HexToHash(utils.IndexTwo)
	LiveMetadata = utils.GetStorageValueMetadata(Live, nil, utils.Uint256)

	VatKey      = common.HexToHash(utils.IndexThree)
	VatMetadata = utils.GetStorageValueMetadata(Vat, nil, utils.Address)

	VowKey      = common.HexToHash(utils.IndexFour)
	VowMetadata = utils.GetStorageValueMetadata(Vow, nil, utils.Address)
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

func (loader *keysLoader) LoadMappings() (map[common.Hash]utils.StorageValueMetadata, error) {
	mappings := loadStaticMappings()
	return loader.addIlkKeys(mappings)
}

func (loader *keysLoader) addIlkKeys(mappings map[common.Hash]utils.StorageValueMetadata) (map[common.Hash]utils.StorageValueMetadata, error) {
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

func loadStaticMappings() map[common.Hash]utils.StorageValueMetadata {
	mappings := make(map[common.Hash]utils.StorageValueMetadata)
	mappings[LiveKey] = LiveMetadata
	mappings[VatKey] = VatMetadata
	mappings[VowKey] = VowMetadata
	return mappings
}

func getIlkFlipKey(ilk string) common.Hash {
	return utils.GetStorageKeyForMapping(IlksMappingIndex, ilk)
}

func getIlkFlipMetadata(ilk string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Ilk: ilk}
	return utils.GetStorageValueMetadata(IlkFlip, keys, utils.Address)
}

func getIlkChopKey(ilk string) common.Hash {
	return utils.GetIncrementedStorageKey(getIlkFlipKey(ilk), 1)
}

func getIlkChopMetadata(ilk string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Ilk: ilk}
	return utils.GetStorageValueMetadata(IlkChop, keys, utils.Uint256)
}

func getIlkLumpKey(ilk string) common.Hash {
	return utils.GetIncrementedStorageKey(getIlkFlipKey(ilk), 2)
}

func getIlkLumpMetadata(ilk string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Ilk: ilk}
	return utils.GetStorageValueMetadata(IlkLump, keys, utils.Uint256)
}
