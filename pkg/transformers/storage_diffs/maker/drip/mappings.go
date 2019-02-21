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
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/maker"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/shared"
)

const (
	IlkTax   = "tax"
	IlkRho   = "rho"
	DripVat  = "vat"
	DripVow  = "vow"
	DripRepo = "repo"
)

var (
	IlkMappingIndex = storage_diffs.IndexOne

	VatKey      = common.HexToHash(storage_diffs.IndexTwo)
	VatMetadata = shared.GetStorageValueMetadata(DripVat, nil, shared.Address)

	VowKey      = common.HexToHash(storage_diffs.IndexThree)
	VowMetadata = shared.GetStorageValueMetadata(DripVow, nil, shared.Bytes32)

	RepoKey      = common.HexToHash(storage_diffs.IndexFour)
	RepoMetadata = shared.GetStorageValueMetadata(DripRepo, nil, shared.Uint256)
)

type DripMappings struct {
	StorageRepository maker.IMakerStorageRepository
	mappings          map[common.Hash]shared.StorageValueMetadata
}

func (mappings *DripMappings) SetDB(db *postgres.DB) {
	mappings.StorageRepository.SetDB(db)
}

func (mappings *DripMappings) Lookup(key common.Hash) (shared.StorageValueMetadata, error) {
	metadata, ok := mappings.mappings[key]
	if !ok {
		err := mappings.loadMappings()
		if err != nil {
			return metadata, err
		}
		metadata, ok = mappings.mappings[key]
		if !ok {
			return metadata, shared.ErrStorageKeyNotFound{Key: key.Hex()}
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

func getStaticMappings() map[common.Hash]shared.StorageValueMetadata {
	mappings := make(map[common.Hash]shared.StorageValueMetadata)
	mappings[VatKey] = VatMetadata
	mappings[VowKey] = VowMetadata
	mappings[RepoKey] = RepoMetadata
	return mappings
}

func getTaxKey(ilk string) common.Hash {
	return storage_diffs.GetMapping(IlkMappingIndex, ilk)
}

func getTaxMetadata(ilk string) shared.StorageValueMetadata {
	keys := map[shared.Key]string{shared.Ilk: ilk}
	return shared.GetStorageValueMetadata(IlkTax, keys, shared.Uint256)
}

func getRhoKey(ilk string) common.Hash {
	return storage_diffs.GetIncrementedKey(getTaxKey(ilk), 1)
}

func getRhoMetadata(ilk string) shared.StorageValueMetadata {
	keys := map[shared.Key]string{shared.Ilk: ilk}
	return shared.GetStorageValueMetadata(IlkRho, keys, shared.Uint48)
}
