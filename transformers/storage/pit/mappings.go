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

package pit

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/vulcanize/vulcanizedb/libraries/shared/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	s2 "github.com/vulcanize/mcd_transformers/transformers/storage"
)

const (
	IlkSpot = "spot"
	PitDrip = "drip"
	PitLine = "Line"
	PitLive = "live"
	PitVat  = "vat"
)

var (
	// storage key and value metadata for "drip" on the Pit contract
	DripKey      = common.HexToHash(storage.IndexFive)
	DripMetadata = utils.StorageValueMetadata{
		Name: PitDrip,
		Keys: nil,
		Type: utils.Address,
	}

	IlkSpotIndex = storage.IndexOne

	// storage key and value metadata for "Spot" on the Pit contract
	LineKey      = common.HexToHash(storage.IndexThree)
	LineMetadata = utils.StorageValueMetadata{
		Name: PitLine,
		Keys: nil,
		Type: utils.Uint256,
	}

	// storage key and value metadata for "live" on the Pit contract
	LiveKey      = common.HexToHash(storage.IndexTwo)
	LiveMetadata = utils.StorageValueMetadata{
		Name: PitLive,
		Keys: nil,
		Type: utils.Uint256,
	}

	// storage key and value metadata for "vat" on the Pit contract
	VatKey      = common.HexToHash(storage.IndexFour)
	VatMetadata = utils.StorageValueMetadata{
		Name: PitVat,
		Keys: nil,
		Type: utils.Address,
	}
)

type PitMappings struct {
	StorageRepository s2.IMakerStorageRepository
	mappings          map[common.Hash]utils.StorageValueMetadata
}

func (mappings *PitMappings) SetDB(db *postgres.DB) {
	mappings.StorageRepository.SetDB(db)
}

func (mappings *PitMappings) Lookup(key common.Hash) (utils.StorageValueMetadata, error) {
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

func (mappings *PitMappings) loadMappings() error {
	mappings.mappings = getStaticMappings()
	ilks, err := mappings.StorageRepository.GetIlks()
	if err != nil {
		return err
	}
	for _, ilk := range ilks {
		mappings.mappings[getSpotKey(ilk)] = getSpotMetadata(ilk)
	}
	return nil
}

func getStaticMappings() map[common.Hash]utils.StorageValueMetadata {
	mappings := make(map[common.Hash]utils.StorageValueMetadata)
	mappings[DripKey] = DripMetadata
	mappings[LineKey] = LineMetadata
	mappings[LiveKey] = LiveMetadata
	mappings[VatKey] = VatMetadata
	return mappings
}

//TODO: remove when Urn query is updated
func getSpotKey(ilk string) common.Hash {
	return storage.GetMapping(IlkSpotIndex, ilk)
}

func getSpotMetadata(ilk string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Ilk: ilk}
	return utils.GetStorageValueMetadata(IlkSpot, keys, utils.Uint256)
}
