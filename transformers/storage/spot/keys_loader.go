// VulcanizeDB
// Copyright © 2019 Vulcanize

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

package spot

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	vdbStorage "github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

const (
	IlkPip = "pip"
	IlkMat = "mat"
	Vat    = "vat"
	Par    = "par"
)

var (
	IlkMappingIndex = vdbStorage.IndexOne

	VatKey      = common.HexToHash(vdbStorage.IndexTwo)
	VatMetadata = vdbStorage.GetStorageValueMetadata(Vat, nil, vdbStorage.Address)

	ParKey      = common.HexToHash(vdbStorage.IndexThree)
	ParMetadata = vdbStorage.GetStorageValueMetadata(Par, nil, vdbStorage.Uint256)
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
	mappings := getStaticMappings()
	ilks, err := loader.storageRepository.GetIlks()
	if err != nil {
		return nil, err
	}
	for _, ilk := range ilks {
		mappings[getPipKey(ilk)] = getPipMetadata(ilk)
		mappings[getMatKey(ilk)] = getMatMetadata(ilk)
	}
	return mappings, nil
}

func getStaticMappings() map[common.Hash]vdbStorage.StorageValueMetadata {
	mappings := make(map[common.Hash]vdbStorage.StorageValueMetadata)
	mappings[VatKey] = VatMetadata
	mappings[ParKey] = ParMetadata
	return mappings
}

func getPipKey(ilk string) common.Hash {
	return vdbStorage.GetStorageKeyForMapping(IlkMappingIndex, ilk)
}

func getPipMetadata(ilk string) vdbStorage.StorageValueMetadata {
	keys := map[vdbStorage.Key]string{constants.Ilk: ilk}
	return vdbStorage.GetStorageValueMetadata(IlkPip, keys, vdbStorage.Address)
}

func getMatKey(ilk string) common.Hash {
	return vdbStorage.GetIncrementedStorageKey(getPipKey(ilk), 1)
}

func getMatMetadata(ilk string) vdbStorage.StorageValueMetadata {
	keys := map[vdbStorage.Key]string{constants.Ilk: ilk}
	return vdbStorage.GetStorageValueMetadata(IlkMat, keys, vdbStorage.Uint256)
}
