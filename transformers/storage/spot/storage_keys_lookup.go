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
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage"
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
	IlkMappingIndex = storage.IndexOne

	VatKey      = common.HexToHash(storage.IndexTwo)
	VatMetadata = utils.GetStorageValueMetadata(Vat, nil, utils.Address)

	ParKey      = common.HexToHash(storage.IndexThree)
	ParMetadata = utils.GetStorageValueMetadata(Par, nil, utils.Uint256)
)

type StorageKeysLookup struct {
	StorageRepository mcdStorage.IMakerStorageRepository
	mappings          map[common.Hash]utils.StorageValueMetadata
}

func (lookup *StorageKeysLookup) SetDB(db *postgres.DB) {
	lookup.StorageRepository.SetDB(db)
}

func (lookup *StorageKeysLookup) Lookup(key common.Hash) (utils.StorageValueMetadata, error) {
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

func (lookup *StorageKeysLookup) loadMappings() error {
	lookup.mappings = getStaticMappings()
	ilks, err := lookup.StorageRepository.GetIlks()
	if err != nil {
		return err
	}
	for _, ilk := range ilks {
		lookup.mappings[getPipKey(ilk)] = getPipMetadata(ilk)
		lookup.mappings[getMatKey(ilk)] = getMatMetadata(ilk)
	}
	lookup.mappings = storage.AddHashedKeys(lookup.mappings)
	return nil
}

func getStaticMappings() map[common.Hash]utils.StorageValueMetadata {
	mappings := make(map[common.Hash]utils.StorageValueMetadata)
	mappings[VatKey] = VatMetadata
	mappings[ParKey] = ParMetadata
	return mappings
}

func getPipKey(ilk string) common.Hash {
	return storage.GetMapping(IlkMappingIndex, ilk)
}

func getPipMetadata(ilk string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Ilk: ilk}
	return utils.GetStorageValueMetadata(IlkPip, keys, utils.Address)
}

func getMatKey(ilk string) common.Hash {
	return storage.GetIncrementedKey(getPipKey(ilk), 1)
}

func getMatMetadata(ilk string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Ilk: ilk}
	return utils.GetStorageValueMetadata(IlkMat, keys, utils.Uint256)
}
