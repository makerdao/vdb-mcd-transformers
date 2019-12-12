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
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	vdbStorage "github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

const (
	IlkDuty = "duty"
	IlkRho  = "rho"
	Vat     = "vat"
	Vow     = "vow"
	Base    = "base"
)

var (
	IlkMappingIndex = vdbStorage.IndexOne

	VatKey      = common.HexToHash(vdbStorage.IndexTwo)
	VatMetadata = vdbStorage.GetValueMetadata(Vat, nil, vdbStorage.Address)

	VowKey      = common.HexToHash(vdbStorage.IndexThree)
	VowMetadata = vdbStorage.GetValueMetadata(Vow, nil, vdbStorage.Bytes32)

	BaseKey      = common.HexToHash(vdbStorage.IndexFour)
	BaseMetadata = vdbStorage.GetValueMetadata(Base, nil, vdbStorage.Uint256)
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

func (loader *keysLoader) LoadMappings() (map[common.Hash]vdbStorage.ValueMetadata, error) {
	mappings := getStaticMappings()
	ilks, err := loader.storageRepository.GetIlks()
	if err != nil {
		return nil, err
	}
	for _, ilk := range ilks {
		mappings[getDutyKey(ilk)] = getDutyMetadata(ilk)
		mappings[getRhoKey(ilk)] = getRhoMetadata(ilk)
	}
	return mappings, nil
}

func getStaticMappings() map[common.Hash]vdbStorage.ValueMetadata {
	mappings := make(map[common.Hash]vdbStorage.ValueMetadata)
	mappings[VatKey] = VatMetadata
	mappings[VowKey] = VowMetadata
	mappings[BaseKey] = BaseMetadata
	return mappings
}

func getDutyKey(ilk string) common.Hash {
	return vdbStorage.GetKeyForMapping(IlkMappingIndex, ilk)
}

func getDutyMetadata(ilk string) vdbStorage.ValueMetadata {
	keys := map[vdbStorage.Key]string{constants.Ilk: ilk}
	return vdbStorage.GetValueMetadata(IlkDuty, keys, vdbStorage.Uint256)
}

func getRhoKey(ilk string) common.Hash {
	return vdbStorage.GetIncrementedKey(getDutyKey(ilk), 1)
}

func getRhoMetadata(ilk string) vdbStorage.ValueMetadata {
	keys := map[vdbStorage.Key]string{constants.Ilk: ilk}
	return vdbStorage.GetValueMetadata(IlkRho, keys, vdbStorage.Uint256)
}
