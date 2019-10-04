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
	s2 "github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/mcd_transformers/transformers/storage/utilities"
)

const (
	Dai     = "dai"
	Gem     = "gem"
	IlkArt  = "Art"
	IlkDust = "dust"
	IlkLine = "line"
	IlkRate = "rate"
	IlkSpot = "spot"
	Sin     = "sin"
	UrnArt  = "art"
	UrnInk  = "ink"
	VatDebt = "debt"
	VatVice = "vice"
	VatLine = "Line"
	VatLive = "live"
)

var (
	IlksMappingIndex = storage.IndexTwo
	UrnsMappingIndex = storage.IndexThree
	GemsMappingIndex = storage.IndexFour
	DaiMappingIndex  = storage.IndexFive
	SinMappingIndex  = storage.IndexSix

	DebtKey      = common.HexToHash(storage.IndexSeven)
	DebtMetadata = utils.StorageValueMetadata{
		Name: VatDebt,
		Keys: nil,
		Type: utils.Uint256,
	}

	ViceKey      = common.HexToHash(storage.IndexEight)
	ViceMetadata = utils.StorageValueMetadata{
		Name: VatVice,
		Keys: nil,
		Type: utils.Uint256,
	}

	LineKey      = common.HexToHash(storage.IndexNine)
	LineMetadata = utils.StorageValueMetadata{
		Name: VatLine,
		Keys: nil,
		Type: utils.Uint256,
	}

	LiveKey      = common.HexToHash(storage.IndexTen)
	LiveMetadata = utils.StorageValueMetadata{
		Name: VatLive,
		Keys: nil,
		Type: utils.Uint256,
	}
)

type VatMappings struct {
	StorageRepository s2.IMakerStorageRepository
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
	mappings.mappings = storage.AddHashedKeys(mappings.mappings)
	return nil
}

func loadStaticMappings() map[common.Hash]utils.StorageValueMetadata {
	mappings := make(map[common.Hash]utils.StorageValueMetadata)
	mappings[DebtKey] = DebtMetadata
	mappings[ViceKey] = ViceMetadata
	mappings[LineKey] = LineMetadata
	mappings[LiveKey] = LiveMetadata
	return mappings
}

func (mappings *VatMappings) loadDaiKeys() error {
	daiKeys, err := mappings.StorageRepository.GetDaiKeys()
	if err != nil {
		return err
	}
	for _, d := range daiKeys {
		paddedDai, padErr := utilities.PadAddress(d)
		if padErr != nil {
			return padErr
		}
		mappings.mappings[getDaiKey(paddedDai)] = getDaiMetadata(d)
	}
	return nil
}

func (mappings *VatMappings) loadGemKeys() error {
	gemKeys, err := mappings.StorageRepository.GetGemKeys()
	if err != nil {
		return err
	}
	for _, gem := range gemKeys {
		paddedGem, padErr := utilities.PadAddress(gem.Identifier)
		if padErr != nil {
			return padErr
		}
		mappings.mappings[getGemKey(gem.Ilk, paddedGem)] = getGemMetadata(gem.Ilk, gem.Identifier)
	}
	return nil
}

func (mappings *VatMappings) loadIlkKeys() error {
	ilks, err := mappings.StorageRepository.GetIlks()
	if err != nil {
		return err
	}
	for _, ilk := range ilks {
		mappings.mappings[getIlkArtKey(ilk)] = getIlkArtMetadata(ilk)
		mappings.mappings[getIlkRateKey(ilk)] = getIlkRateMetadata(ilk)
		mappings.mappings[getIlkSpotKey(ilk)] = getIlkSpotMetadata(ilk)
		mappings.mappings[getIlkLineKey(ilk)] = getIlkLineMetadata(ilk)
		mappings.mappings[getIlkDustKey(ilk)] = getIlkDustMetadata(ilk)
	}
	return nil
}

func (mappings *VatMappings) loadSinKeys() error {
	sinKeys, err := mappings.StorageRepository.GetVatSinKeys()
	if err != nil {
		return err
	}
	for _, s := range sinKeys {
		paddedSin, padErr := utilities.PadAddress(s)
		if padErr != nil {
			return padErr
		}
		mappings.mappings[getSinKey(paddedSin)] = getSinMetadata(s)
	}
	return nil
}

func (mappings *VatMappings) loadUrnKeys() error {
	urns, err := mappings.StorageRepository.GetUrns()
	if err != nil {
		return err
	}
	for _, urn := range urns {
		paddedGuy, padErr := utilities.PadAddress(urn.Identifier)
		if padErr != nil {
			return padErr
		}
		mappings.mappings[getUrnArtKey(urn.Ilk, paddedGuy)] = getUrnArtMetadata(urn.Ilk, urn.Identifier)
		mappings.mappings[getUrnInkKey(urn.Ilk, paddedGuy)] = getUrnInkMetadata(urn.Ilk, urn.Identifier)
	}
	return nil
}

func getIlkArtKey(ilk string) common.Hash {
	return storage.GetMapping(IlksMappingIndex, ilk)
}

func getIlkArtMetadata(ilk string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Ilk: ilk}
	return utils.GetStorageValueMetadata(IlkArt, keys, utils.Uint256)
}

func getIlkRateKey(ilk string) common.Hash {
	return storage.GetIncrementedKey(getIlkArtKey(ilk), 1)
}

func getIlkRateMetadata(ilk string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Ilk: ilk}
	return utils.GetStorageValueMetadata(IlkRate, keys, utils.Uint256)
}

func getIlkSpotKey(ilk string) common.Hash {
	return storage.GetIncrementedKey(getIlkArtKey(ilk), 2)
}

func getIlkSpotMetadata(ilk string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Ilk: ilk}
	return utils.GetStorageValueMetadata(IlkSpot, keys, utils.Uint256)
}

func getIlkLineKey(ilk string) common.Hash {
	return storage.GetIncrementedKey(getIlkArtKey(ilk), 3)
}

func getIlkLineMetadata(ilk string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Ilk: ilk}
	return utils.GetStorageValueMetadata(IlkLine, keys, utils.Uint256)
}

func getIlkDustKey(ilk string) common.Hash {
	return storage.GetIncrementedKey(getIlkArtKey(ilk), 4)
}

func getIlkDustMetadata(ilk string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Ilk: ilk}
	return utils.GetStorageValueMetadata(IlkDust, keys, utils.Uint256)
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
