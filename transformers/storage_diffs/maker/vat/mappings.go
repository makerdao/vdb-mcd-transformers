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

package vat

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/vulcanize/vulcanizedb/libraries/shared/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage_diffs/maker"
)

const (
	Dai     = "dai"
	Gem     = "gem"
	IlkArt  = "Art"
	IlkInk  = "Ink"
	IlkRate = "rate"
	IlkTake = "take"
	Sin     = "sin"
	UrnArt  = "art"
	UrnInk  = "ink"
	VatDebt = "debt"
	VatVice = "vice"
)

var (
	DebtKey      = common.HexToHash(storage.IndexSix)
	DebtMetadata = utils.StorageValueMetadata{
		Name: VatDebt,
		Keys: nil,
		Type: 0,
	}

	IlksMappingIndex = storage.IndexOne
	UrnsMappingIndex = storage.IndexTwo
	GemsMappingIndex = storage.IndexThree
	DaiMappingIndex  = storage.IndexFour
	SinMappingIndex  = storage.IndexFive

	ViceKey      = common.HexToHash(storage.IndexSeven)
	ViceMetadata = utils.StorageValueMetadata{
		Name: VatVice,
		Keys: nil,
		Type: 0,
	}
)

type VatMappings struct {
	StorageRepository maker.IMakerStorageRepository
	mappings          map[common.Hash]utils.StorageValueMetadata
}

func (mappings VatMappings) Lookup(key common.Hash) (utils.StorageValueMetadata, error) {
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

func (mappings *VatMappings) SetDB(db *postgres.DB) {
	mappings.StorageRepository.SetDB(db)
}

func (mappings *VatMappings) loadMappings() error {
	mappings.mappings = loadStaticMappings()
	daiErr := mappings.loadDaiKeys()
	if daiErr != nil {
		return daiErr
	}
	gemErr := mappings.loadGemKeys()
	if gemErr != nil {
		return gemErr
	}
	ilkErr := mappings.loadIlkKeys()
	if ilkErr != nil {
		return ilkErr
	}
	sinErr := mappings.loadSinKeys()
	if sinErr != nil {
		return sinErr
	}
	urnErr := mappings.loadUrnKeys()
	if urnErr != nil {
		return urnErr
	}
	return nil
}

func loadStaticMappings() map[common.Hash]utils.StorageValueMetadata {
	mappings := make(map[common.Hash]utils.StorageValueMetadata)
	mappings[DebtKey] = DebtMetadata
	mappings[ViceKey] = ViceMetadata
	return mappings
}

func (mappings *VatMappings) loadDaiKeys() error {
	daiKeys, err := mappings.StorageRepository.GetDaiKeys()
	if err != nil {
		return err
	}
	for _, d := range daiKeys {
		mappings.mappings[getDaiKey(d)] = getDaiMetadata(d)
	}
	return nil
}

func (mappings *VatMappings) loadGemKeys() error {
	gemKeys, err := mappings.StorageRepository.GetGemKeys()
	if err != nil {
		return err
	}
	for _, gem := range gemKeys {
		mappings.mappings[getGemKey(gem.Ilk, gem.Guy)] = getGemMetadata(gem.Ilk, gem.Guy)
	}
	return nil
}

func (mappings *VatMappings) loadIlkKeys() error {
	ilks, err := mappings.StorageRepository.GetIlks()
	if err != nil {
		return err
	}
	for _, ilk := range ilks {
		mappings.mappings[getIlkTakeKey(ilk)] = getIlkTakeMetadata(ilk)
		mappings.mappings[getIlkRateKey(ilk)] = getIlkRateMetadata(ilk)
		mappings.mappings[getIlkInkKey(ilk)] = getIlkInkMetadata(ilk)
		mappings.mappings[getIlkArtKey(ilk)] = getIlkArtMetadata(ilk)
	}
	return nil
}

func (mappings *VatMappings) loadSinKeys() error {
	sinKeys, err := mappings.StorageRepository.GetSinKeys()
	if err != nil {
		return err
	}
	for _, s := range sinKeys {
		mappings.mappings[getSinKey(s)] = getSinMetadata(s)
	}
	return nil
}

func (mappings *VatMappings) loadUrnKeys() error {
	urns, err := mappings.StorageRepository.GetUrns()
	if err != nil {
		return err
	}
	for _, urn := range urns {
		mappings.mappings[getUrnInkKey(urn.Ilk, urn.Guy)] = getUrnInkMetadata(urn.Ilk, urn.Guy)
		mappings.mappings[getUrnArtKey(urn.Ilk, urn.Guy)] = getUrnArtMetadata(urn.Ilk, urn.Guy)
	}
	return nil
}

func getIlkTakeKey(ilk string) common.Hash {
	return storage.GetMapping(IlksMappingIndex, ilk)
}

func getIlkTakeMetadata(ilk string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Ilk: ilk}
	return utils.GetStorageValueMetadata(IlkTake, keys, utils.Uint256)
}

func getIlkRateKey(ilk string) common.Hash {
	return storage.GetIncrementedKey(getIlkTakeKey(ilk), 1)
}

func getIlkRateMetadata(ilk string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Ilk: ilk}
	return utils.GetStorageValueMetadata(IlkRate, keys, utils.Uint256)
}

func getIlkInkKey(ilk string) common.Hash {
	return storage.GetIncrementedKey(getIlkTakeKey(ilk), 2)
}

func getIlkInkMetadata(ilk string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Ilk: ilk}
	return utils.GetStorageValueMetadata(IlkInk, keys, utils.Uint256)
}

func getIlkArtKey(ilk string) common.Hash {
	return storage.GetIncrementedKey(getIlkTakeKey(ilk), 3)
}

func getIlkArtMetadata(ilk string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Ilk: ilk}
	return utils.GetStorageValueMetadata(IlkArt, keys, utils.Uint256)
}

func getUrnInkKey(ilk, guy string) common.Hash {
	return storage.GetNestedMapping(UrnsMappingIndex, ilk, guy)
}

func getUrnInkMetadata(ilk, guy string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Ilk: ilk, constants.Guy: guy}
	return utils.GetStorageValueMetadata(UrnInk, keys, utils.Uint256)
}

func getUrnArtKey(ilk, guy string) common.Hash {
	return storage.GetIncrementedKey(getUrnInkKey(ilk, guy), 1)
}

func getUrnArtMetadata(ilk, guy string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Ilk: ilk, constants.Guy: guy}
	return utils.GetStorageValueMetadata(UrnArt, keys, utils.Uint256)
}

func getGemKey(ilk, guy string) common.Hash {
	return storage.GetNestedMapping(GemsMappingIndex, ilk, guy)
}

func getGemMetadata(ilk, guy string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Ilk: ilk, constants.Guy: guy}
	return utils.GetStorageValueMetadata(Gem, keys, utils.Uint256)
}

func getDaiKey(guy string) common.Hash {
	return storage.GetMapping(DaiMappingIndex, guy)
}

func getDaiMetadata(guy string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Guy: guy}
	return utils.GetStorageValueMetadata(Dai, keys, utils.Uint256)
}

func getSinKey(guy string) common.Hash {
	return storage.GetMapping(SinMappingIndex, guy)
}

func getSinMetadata(guy string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Guy: guy}
	return utils.GetStorageValueMetadata(Sin, keys, utils.Uint256)
}
