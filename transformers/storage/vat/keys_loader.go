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
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	mcdStorage "github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/mcd_transformers/transformers/storage/utilities"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
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
	Debt    = "debt"
	Vice    = "vice"
	Line    = "Line"
	Live    = "live"
)

var (
	IlksMappingIndex = storage.IndexTwo
	UrnsMappingIndex = storage.IndexThree
	GemsMappingIndex = storage.IndexFour
	DaiMappingIndex  = storage.IndexFive
	SinMappingIndex  = storage.IndexSix

	DebtKey      = common.HexToHash(storage.IndexSeven)
	DebtMetadata = utils.StorageValueMetadata{
		Name: Debt,
		Keys: nil,
		Type: utils.Uint256,
	}

	ViceKey      = common.HexToHash(storage.IndexEight)
	ViceMetadata = utils.StorageValueMetadata{
		Name: Vice,
		Keys: nil,
		Type: utils.Uint256,
	}

	LineKey      = common.HexToHash(storage.IndexNine)
	LineMetadata = utils.StorageValueMetadata{
		Name: Line,
		Keys: nil,
		Type: utils.Uint256,
	}

	LiveKey      = common.HexToHash(storage.IndexTen)
	LiveMetadata = utils.StorageValueMetadata{
		Name: Live,
		Keys: nil,
		Type: utils.Uint256,
	}
)

type keysLoader struct {
	storageRepository mcdStorage.IMakerStorageRepository
}

func NewKeysLoader(storageRepository mcdStorage.IMakerStorageRepository) mcdStorage.KeysLoader {
	return &keysLoader{storageRepository: storageRepository}
}

func (loader *keysLoader) SetDB(db *postgres.DB) {
	loader.storageRepository.SetDB(db)
}

func (loader *keysLoader) LoadMappings() (map[common.Hash]utils.StorageValueMetadata, error) {
	mappings := loadStaticMappings()
	mappings, daiErr := loader.addDaiKeys(mappings)
	if daiErr != nil {
		return nil, daiErr
	}
	mappings, gemErr := loader.addGemKeys(mappings)
	if gemErr != nil {
		return nil, gemErr
	}
	mappings, ilkErr := loader.addIlkKeys(mappings)
	if ilkErr != nil {
		return nil, ilkErr
	}
	mappings, sinErr := loader.addSinKeys(mappings)
	if sinErr != nil {
		return nil, sinErr
	}
	return loader.addUrnKeys(mappings)
}

func loadStaticMappings() map[common.Hash]utils.StorageValueMetadata {
	mappings := make(map[common.Hash]utils.StorageValueMetadata)
	mappings[DebtKey] = DebtMetadata
	mappings[ViceKey] = ViceMetadata
	mappings[LineKey] = LineMetadata
	mappings[LiveKey] = LiveMetadata
	return mappings
}

func (loader *keysLoader) addDaiKeys(mappings map[common.Hash]utils.StorageValueMetadata) (map[common.Hash]utils.StorageValueMetadata, error) {
	daiKeys, err := loader.storageRepository.GetDaiKeys()
	if err != nil {
		return nil, err
	}
	for _, d := range daiKeys {
		paddedDai, padErr := utilities.PadAddress(d)
		if padErr != nil {
			return nil, padErr
		}
		mappings[getDaiKey(paddedDai)] = getDaiMetadata(d)
	}
	return mappings, nil
}

func (loader *keysLoader) addGemKeys(mappings map[common.Hash]utils.StorageValueMetadata) (map[common.Hash]utils.StorageValueMetadata, error) {
	gemKeys, err := loader.storageRepository.GetGemKeys()
	if err != nil {
		return nil, err
	}
	for _, gem := range gemKeys {
		paddedGem, padErr := utilities.PadAddress(gem.Identifier)
		if padErr != nil {
			return nil, padErr
		}
		mappings[getGemKey(gem.Ilk, paddedGem)] = getGemMetadata(gem.Ilk, gem.Identifier)
	}
	return mappings, nil
}

func (loader *keysLoader) addIlkKeys(mappings map[common.Hash]utils.StorageValueMetadata) (map[common.Hash]utils.StorageValueMetadata, error) {
	ilks, err := loader.storageRepository.GetIlks()
	if err != nil {
		return nil, err
	}
	for _, ilk := range ilks {
		mappings[getIlkArtKey(ilk)] = getIlkArtMetadata(ilk)
		mappings[getIlkRateKey(ilk)] = getIlkRateMetadata(ilk)
		mappings[getIlkSpotKey(ilk)] = getIlkSpotMetadata(ilk)
		mappings[getIlkLineKey(ilk)] = getIlkLineMetadata(ilk)
		mappings[getIlkDustKey(ilk)] = getIlkDustMetadata(ilk)
	}
	return mappings, nil
}

func (loader *keysLoader) addSinKeys(mappings map[common.Hash]utils.StorageValueMetadata) (map[common.Hash]utils.StorageValueMetadata, error) {
	sinKeys, err := loader.storageRepository.GetVatSinKeys()
	if err != nil {
		return nil, err
	}
	for _, s := range sinKeys {
		paddedSin, padErr := utilities.PadAddress(s)
		if padErr != nil {
			return nil, padErr
		}
		mappings[getSinKey(paddedSin)] = getSinMetadata(s)
	}
	return mappings, nil
}

func (loader *keysLoader) addUrnKeys(mappings map[common.Hash]utils.StorageValueMetadata) (map[common.Hash]utils.StorageValueMetadata, error) {
	urns, err := loader.storageRepository.GetUrns()
	if err != nil {
		return nil, err
	}
	for _, urn := range urns {
		paddedGuy, padErr := utilities.PadAddress(urn.Identifier)
		if padErr != nil {
			return nil, padErr
		}
		mappings[getUrnArtKey(urn.Ilk, paddedGuy)] = getUrnArtMetadata(urn.Ilk, urn.Identifier)
		mappings[getUrnInkKey(urn.Ilk, paddedGuy)] = getUrnInkMetadata(urn.Ilk, urn.Identifier)
	}
	return mappings, nil
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
