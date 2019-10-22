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
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	mcdStorage "github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/factories/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

const (
	IlkPip = "pip"
	IlkMat = "mat"
	Vat    = "vat"
	Par    = "par"
)

var (
	IlkMappingIndex = utils.IndexOne

	VatKey      = common.HexToHash(utils.IndexTwo)
	VatMetadata = utils.GetStorageValueMetadata(Vat, nil, utils.Address)

	ParKey      = common.HexToHash(utils.IndexThree)
	ParMetadata = utils.GetStorageValueMetadata(Par, nil, utils.Uint256)
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

func getStaticMappings() map[common.Hash]utils.StorageValueMetadata {
	mappings := make(map[common.Hash]utils.StorageValueMetadata)
	mappings[VatKey] = VatMetadata
	mappings[ParKey] = ParMetadata
	return mappings
}

func getPipKey(ilk string) common.Hash {
	return utils.GetStorageKeyForMapping(IlkMappingIndex, ilk)
}

func getPipMetadata(ilk string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Ilk: ilk}
	return utils.GetStorageValueMetadata(IlkPip, keys, utils.Address)
}

func getMatKey(ilk string) common.Hash {
	return utils.GetIncrementedStorageKey(getPipKey(ilk), 1)
}

func getMatMetadata(ilk string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Ilk: ilk}
	return utils.GetStorageValueMetadata(IlkMat, keys, utils.Uint256)
}
