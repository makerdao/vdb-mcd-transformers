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

package drip

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage_diffs/maker"

	"github.com/vulcanize/vulcanizedb/libraries/shared/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

const (
	IlkTax   = "tax"
	IlkRho   = "rho"
	DripVat  = "vat"
	DripVow  = "vow"
	DripRepo = "repo"
)

var (
	IlkMappingIndex = storage.IndexOne

	VatKey      = common.HexToHash(storage.IndexTwo)
	VatMetadata = utils.GetStorageValueMetadata(DripVat, nil, utils.Address)

	VowKey      = common.HexToHash(storage.IndexThree)
	VowMetadata = utils.GetStorageValueMetadata(DripVow, nil, utils.Bytes32)

	RepoKey      = common.HexToHash(storage.IndexFour)
	RepoMetadata = utils.GetStorageValueMetadata(DripRepo, nil, utils.Uint256)
)

type DripMappings struct {
	StorageRepository maker.IMakerStorageRepository
	mappings          map[common.Hash]utils.StorageValueMetadata
}

func (mappings *DripMappings) SetDB(db *postgres.DB) {
	mappings.StorageRepository.SetDB(db)
}

func (mappings *DripMappings) Lookup(key common.Hash) (utils.StorageValueMetadata, error) {
	metadata, ok := mappings.mappings[key]
	if !ok {
		err := mappings.loadMappings()
		if err != nil {
			return metadata, err
		}
		metadata, ok = mappings.mappings[key]
		if !ok {
			return metadata, utils.ErrStorageKeyNotFound{Key: key.Hex()}
		}
	}
	return metadata, nil
}

func (mappings *DripMappings) loadMappings() error {
	mappings.mappings = getStaticMappings()
	ilks, err := mappings.StorageRepository.GetIlks()
	if err != nil {
		return err
	}
	for _, ilk := range ilks {
		mappings.mappings[getTaxKey(ilk)] = getTaxMetadata(ilk)
		mappings.mappings[getRhoKey(ilk)] = getRhoMetadata(ilk)
	}
	return nil
}

func getStaticMappings() map[common.Hash]utils.StorageValueMetadata {
	mappings := make(map[common.Hash]utils.StorageValueMetadata)
	mappings[VatKey] = VatMetadata
	mappings[VowKey] = VowMetadata
	mappings[RepoKey] = RepoMetadata
	return mappings
}

func getTaxKey(ilk string) common.Hash {
	return storage.GetMapping(IlkMappingIndex, ilk)
}

func getTaxMetadata(ilk string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Ilk: ilk}
	return utils.GetStorageValueMetadata(IlkTax, keys, utils.Uint256)
}

func getRhoKey(ilk string) common.Hash {
	return storage.GetIncrementedKey(getTaxKey(ilk), 1)
}

func getRhoMetadata(ilk string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Ilk: ilk}
	return utils.GetStorageValueMetadata(IlkRho, keys, utils.Uint48)
}
