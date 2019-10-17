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
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	mcdStorage "github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

const (
	VowVat     = "vat"
	VowFlapper = "flapper"
	VowFlopper = "flopper"
	SinMapping = "sin"
	SinInteger = "Sin"
	VowAsh     = "Ash"
	VowWait    = "wait"
	VowDump    = "dump"
	VowSump    = "sump"
	VowBump    = "bump"
	VowHump    = "hump"
)

var (
	VatKey      = common.HexToHash(storage.IndexOne)
	VatMetadata = utils.StorageValueMetadata{
		Name: VowVat,
		Keys: nil,
		Type: utils.Address,
	}

	FlapperKey      = common.HexToHash(storage.IndexTwo)
	FlapperMetadata = utils.StorageValueMetadata{
		Name: VowFlapper,
		Keys: nil,
		Type: utils.Address,
	}

	FlopperKey      = common.HexToHash(storage.IndexThree)
	FlopperMetadata = utils.StorageValueMetadata{
		Name: VowFlopper,
		Keys: nil,
		Type: utils.Address,
	}

	SinMappingIndex = storage.IndexFour

	SinIntegerKey      = common.HexToHash(storage.IndexFive)
	SinIntegerMetadata = utils.StorageValueMetadata{
		Name: SinInteger,
		Keys: nil,
		Type: utils.Uint256,
	}

	AshKey      = common.HexToHash(storage.IndexSix)
	AshMetadata = utils.StorageValueMetadata{
		Name: VowAsh,
		Keys: nil,
		Type: utils.Uint256,
	}

	WaitKey      = common.HexToHash(storage.IndexSeven)
	WaitMetadata = utils.StorageValueMetadata{
		Name: VowWait,
		Keys: nil,
		Type: utils.Uint256,
	}

	DumpKey      = common.HexToHash(storage.IndexEight)
	DumpMetadata = utils.StorageValueMetadata{
		Name: VowDump,
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

type StorageKeysLookup struct {
	StorageRepository mcdStorage.IMakerStorageRepository
	mappings          map[common.Hash]utils.StorageValueMetadata
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
	lookup.mappings = loadStaticMappings()
	sinErr := lookup.loadSinKeys()
	if sinErr != nil {
		return sinErr
	}
	lookup.mappings = storage.AddHashedKeys(lookup.mappings)
	return nil
}

func (lookup *StorageKeysLookup) SetDB(db *postgres.DB) {
	lookup.StorageRepository.SetDB(db)
}

func loadStaticMappings() map[common.Hash]utils.StorageValueMetadata {
	mappings := make(map[common.Hash]utils.StorageValueMetadata)
	mappings[VatKey] = VatMetadata
	mappings[FlapperKey] = FlapperMetadata
	mappings[FlopperKey] = FlopperMetadata
	mappings[SinIntegerKey] = SinIntegerMetadata
	mappings[AshKey] = AshMetadata
	mappings[WaitKey] = WaitMetadata
	mappings[DumpKey] = DumpMetadata
	mappings[SumpKey] = SumpMetadata
	mappings[BumpKey] = BumpMetadata
	mappings[HumpKey] = HumpMetadata
	return mappings
}

func (lookup *StorageKeysLookup) loadSinKeys() error {
	sinKeys, err := lookup.StorageRepository.GetVowSinKeys()
	if err != nil {
		return err
	}
	for _, timestamp := range sinKeys {
		hexTimestamp, err := shared.ConvertIntStringToHex(timestamp)
		if err != nil {
			return err
		}
		lookup.mappings[getSinKey(hexTimestamp)] = getSinMetadata(timestamp)
	}
	return nil
}

func getSinKey(hexTimestamp string) common.Hash {
	return storage.GetMapping(SinMappingIndex, hexTimestamp)
}

func getSinMetadata(timestamp string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Timestamp: timestamp}
	return utils.GetStorageValueMetadata(SinMapping, keys, utils.Uint256)
}
