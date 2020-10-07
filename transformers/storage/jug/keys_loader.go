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
	IlkDuty = "duty"
	IlkRho  = "rho"
	Vat     = "vat"
	Vow     = "vow"
	Base    = "base"
)

var (
	IlkMappingIndex = vdbStorage.IndexOne

	VatKey      = common.HexToHash(vdbStorage.IndexTwo)
	VatMetadata = types.GetValueMetadata(Vat, nil, types.Address)

	VowKey      = common.HexToHash(vdbStorage.IndexThree)
	VowMetadata = types.GetValueMetadata(Vow, nil, types.Bytes32)

	BaseKey      = common.HexToHash(vdbStorage.IndexFour)
	BaseMetadata = types.GetValueMetadata(Base, nil, types.Uint256)
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
	mappings := addStaticMappings()
	mappings, wardsErr := loader.addWardsKeys(mappings)
	if wardsErr != nil {
		return nil, fmt.Errorf("error adding wards keys to jug keys loader: %w", wardsErr)
	}
	mappings, ilksErr := loader.addIlksKeys(mappings)
	if ilksErr != nil {
		return nil, fmt.Errorf("error adding ilks keys to jug keys loader: %w", ilksErr)
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

func (loader *keysLoader) addIlksKeys(mappings map[common.Hash]types.ValueMetadata) (map[common.Hash]types.ValueMetadata, error) {
	ilks, err := loader.storageRepository.GetIlks()
	if err != nil {
		return nil, fmt.Errorf("error getting ilks: %w", err)
	}
	for _, ilk := range ilks {
		mappings[getDutyKey(ilk)] = getDutyMetadata(ilk)
		mappings[getRhoKey(ilk)] = getRhoMetadata(ilk)
	}
	return mappings, nil
}

func addStaticMappings() map[common.Hash]types.ValueMetadata {
	mappings := make(map[common.Hash]types.ValueMetadata)
	mappings[VatKey] = VatMetadata
	mappings[VowKey] = VowMetadata
	mappings[BaseKey] = BaseMetadata
	return mappings
}

func getDutyKey(ilk string) common.Hash {
	return vdbStorage.GetKeyForMapping(IlkMappingIndex, ilk)
}

func getDutyMetadata(ilk string) types.ValueMetadata {
	keys := map[types.Key]string{constants.Ilk: ilk}
	return types.GetValueMetadata(IlkDuty, keys, types.Uint256)
}

func getRhoKey(ilk string) common.Hash {
	return vdbStorage.GetIncrementedKey(getDutyKey(ilk), 1)
}

func getRhoMetadata(ilk string) types.ValueMetadata {
	keys := map[types.Key]string{constants.Ilk: ilk}
	return types.GetValueMetadata(IlkRho, keys, types.Uint256)
}
