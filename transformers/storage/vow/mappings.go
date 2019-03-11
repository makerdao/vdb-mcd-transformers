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

package vow

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/vulcanize/vulcanizedb/libraries/shared/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	s2 "github.com/vulcanize/mcd_transformers/transformers/storage"
)

const (
	VowVat  = "vat"
	VowCow  = "cow"
	VowRow  = "row"
	VowSin  = "Sin"
	VowWoe  = "Woe"
	VowAsh  = "Ash"
	VowWait = "wait"
	VowSump = "sump"
	VowBump = "bump"
	VowHump = "hump"
)

var (
	VatKey      = common.HexToHash(storage.IndexOne)
	VatMetadata = utils.StorageValueMetadata{
		Name: VowVat,
		Keys: nil,
		Type: utils.Address,
	}

	CowKey      = common.HexToHash(storage.IndexTwo)
	CowMetadata = utils.StorageValueMetadata{
		Name: VowCow,
		Keys: nil,
		Type: utils.Address,
	}

	RowKey      = common.HexToHash(storage.IndexThree)
	RowMetadata = utils.StorageValueMetadata{
		Name: VowRow,
		Keys: nil,
		Type: utils.Address,
	}

	SinKey      = common.HexToHash(storage.IndexFive)
	SinMetadata = utils.StorageValueMetadata{
		Name: VowSin,
		Keys: nil,
		Type: utils.Uint256,
	}

	WoeKey      = common.HexToHash(storage.IndexSix)
	WoeMetadata = utils.StorageValueMetadata{
		Name: VowWoe,
		Keys: nil,
		Type: utils.Uint256,
	}

	AshKey      = common.HexToHash(storage.IndexSeven)
	AshMetadata = utils.StorageValueMetadata{
		Name: VowAsh,
		Keys: nil,
		Type: utils.Uint256,
	}

	WaitKey      = common.HexToHash(storage.IndexEight)
	WaitMetadata = utils.StorageValueMetadata{
		Name: VowWait,
		Keys: nil,
		Type: utils.Uint256,
	}

	SumpKey      = common.HexToHash(storage.IndexNine)
	SumpMetadata = utils.StorageValueMetadata{
		Name: VowSump,
		Keys: nil,
		Type: utils.Uint256,
	}

	BumpKey      = common.HexToHash(storage.IndexTen)
	BumpMetadata = utils.StorageValueMetadata{
		Name: VowBump,
		Keys: nil,
		Type: utils.Uint256,
	}

	HumpKey      = common.HexToHash(storage.IndexEleven)
	HumpMetadata = utils.StorageValueMetadata{
		Name: VowHump,
		Keys: nil,
		Type: utils.Uint256,
	}
)

type VowMappings struct {
	StorageRepository s2.IMakerStorageRepository
	mappings          map[common.Hash]utils.StorageValueMetadata
}

func (mappings *VowMappings) Lookup(key common.Hash) (utils.StorageValueMetadata, error) {
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

func (mappings *VowMappings) loadMappings() error {
	staticMappings := make(map[common.Hash]utils.StorageValueMetadata)
	staticMappings[VatKey] = VatMetadata
	staticMappings[CowKey] = CowMetadata
	staticMappings[RowKey] = RowMetadata
	staticMappings[SinKey] = SinMetadata
	staticMappings[WoeKey] = WoeMetadata
	staticMappings[AshKey] = AshMetadata
	staticMappings[WaitKey] = WaitMetadata
	staticMappings[SumpKey] = SumpMetadata
	staticMappings[BumpKey] = BumpMetadata
	staticMappings[HumpKey] = HumpMetadata

	mappings.mappings = staticMappings

	return nil
}

func (mappings *VowMappings) SetDB(db *postgres.DB) {
	mappings.StorageRepository.SetDB(db)
}
