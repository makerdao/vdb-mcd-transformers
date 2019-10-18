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

package jug

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	mcdStorage "github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

const (
	IlkDuty = "duty"
	IlkRho  = "rho"
	Vat     = "vat"
	Vow     = "vow"
	Base    = "base"
)

var (
	IlkMappingIndex = storage.IndexOne

	VatKey      = common.HexToHash(storage.IndexTwo)
	VatMetadata = utils.GetStorageValueMetadata(Vat, nil, utils.Address)

	VowKey      = common.HexToHash(storage.IndexThree)
	VowMetadata = utils.GetStorageValueMetadata(Vow, nil, utils.Bytes32)

	BaseKey      = common.HexToHash(storage.IndexFour)
	BaseMetadata = utils.GetStorageValueMetadata(Base, nil, utils.Uint256)
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
		lookup.mappings[getDutyKey(ilk)] = getDutyMetadata(ilk)
		lookup.mappings[getRhoKey(ilk)] = getRhoMetadata(ilk)
	}
	lookup.mappings = storage.AddHashedKeys(lookup.mappings)
	return nil
}

func getStaticMappings() map[common.Hash]utils.StorageValueMetadata {
	mappings := make(map[common.Hash]utils.StorageValueMetadata)
	mappings[VatKey] = VatMetadata
	mappings[VowKey] = VowMetadata
	mappings[BaseKey] = BaseMetadata
	return mappings
}

func getDutyKey(ilk string) common.Hash {
	return storage.GetMapping(IlkMappingIndex, ilk)
}

func getDutyMetadata(ilk string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Ilk: ilk}
	return utils.GetStorageValueMetadata(IlkDuty, keys, utils.Uint256)
}

func getRhoKey(ilk string) common.Hash {
	return storage.GetIncrementedKey(getDutyKey(ilk), 1)
}

func getRhoMetadata(ilk string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Ilk: ilk}
	return utils.GetStorageValueMetadata(IlkRho, keys, utils.Uint256)
}
