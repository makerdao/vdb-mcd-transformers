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
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities/wards"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	vdbStorage "github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

const (
	Can     = "can"
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
	CanMappingIndex  = vdbStorage.IndexOne
	IlksMappingIndex = vdbStorage.IndexTwo
	UrnsMappingIndex = vdbStorage.IndexThree
	GemsMappingIndex = vdbStorage.IndexFour
	DaiMappingIndex  = vdbStorage.IndexFive
	SinMappingIndex  = vdbStorage.IndexSix

	DebtKey      = common.HexToHash(vdbStorage.IndexSeven)
	DebtMetadata = types.ValueMetadata{
		Name: Debt,
		Keys: nil,
		Type: types.Uint256,
	}

	ViceKey      = common.HexToHash(vdbStorage.IndexEight)
	ViceMetadata = types.ValueMetadata{
		Name: Vice,
		Keys: nil,
		Type: types.Uint256,
	}

	LineKey      = common.HexToHash(vdbStorage.IndexNine)
	LineMetadata = types.ValueMetadata{
		Name: Line,
		Keys: nil,
		Type: types.Uint256,
	}

	LiveKey      = common.HexToHash(vdbStorage.IndexTen)
	LiveMetadata = types.ValueMetadata{
		Name: Live,
		Keys: nil,
		Type: types.Uint256,
	}
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

func (loader *keysLoader) LoadMappings() (map[common.Hash]types.ValueMetadata, error) {
	mappings := loadStaticMappings()
	mappings, wardsErr := loader.addWardsKeys(mappings)
	if wardsErr != nil {
		return nil, wardsErr
	}
	mappings, canErr := loader.addCanKeys(mappings)
	if canErr != nil {
		return nil, canErr
	}
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

func loadStaticMappings() map[common.Hash]types.ValueMetadata {
	mappings := make(map[common.Hash]types.ValueMetadata)
	mappings[DebtKey] = DebtMetadata
	mappings[ViceKey] = ViceMetadata
	mappings[LineKey] = LineMetadata
	mappings[LiveKey] = LiveMetadata
	return mappings
}

func (loader *keysLoader) addCanKeys(mappings map[common.Hash]types.ValueMetadata) (map[common.Hash]types.ValueMetadata, error) {
	cans, err := loader.storageRepository.GetVatCanKeys()
	if err != nil {
		return nil, err
	}
	for _, can := range cans {
		paddedBit, padBitErr := utilities.PadAddress(can.Bit)
		if padBitErr != nil {
			return nil, padBitErr
		}
		paddedUsr, padUsrErr := utilities.PadAddress(can.Usr)
		if padUsrErr != nil {
			return nil, padUsrErr
		}
		mappings[getCanKey(paddedBit, paddedUsr)] = getCanMetadata(can.Bit, can.Usr)
	}
	return mappings, nil
}

func (loader *keysLoader) addWardsKeys(mappings map[common.Hash]types.ValueMetadata) (map[common.Hash]types.ValueMetadata, error) {
	addresses, err := loader.storageRepository.GetVatWardsAddresses()
	if err != nil {
		return nil, err
	}
	return wards.AddWardsKeys(mappings, addresses)
}

func (loader *keysLoader) addDaiKeys(mappings map[common.Hash]types.ValueMetadata) (map[common.Hash]types.ValueMetadata, error) {
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

func (loader *keysLoader) addGemKeys(mappings map[common.Hash]types.ValueMetadata) (map[common.Hash]types.ValueMetadata, error) {
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

func (loader *keysLoader) addIlkKeys(mappings map[common.Hash]types.ValueMetadata) (map[common.Hash]types.ValueMetadata, error) {
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

func (loader *keysLoader) addSinKeys(mappings map[common.Hash]types.ValueMetadata) (map[common.Hash]types.ValueMetadata, error) {
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

func (loader *keysLoader) addUrnKeys(mappings map[common.Hash]types.ValueMetadata) (map[common.Hash]types.ValueMetadata, error) {
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

func getCanKey(bit, usr string) common.Hash {
	return vdbStorage.GetKeyForNestedMapping(CanMappingIndex, bit, usr)
}

func getCanMetadata(bit, usr string) types.ValueMetadata {
	keys := map[types.Key]string{constants.Bit: bit, constants.User: usr}
	return types.GetValueMetadata(Can, keys, types.Uint256)
}

func getIlkArtKey(ilk string) common.Hash {
	return vdbStorage.GetKeyForMapping(IlksMappingIndex, ilk)
}

func getIlkArtMetadata(ilk string) types.ValueMetadata {
	keys := map[types.Key]string{constants.Ilk: ilk}
	return types.GetValueMetadata(IlkArt, keys, types.Uint256)
}

func getIlkRateKey(ilk string) common.Hash {
	return vdbStorage.GetIncrementedKey(getIlkArtKey(ilk), 1)
}

func getIlkRateMetadata(ilk string) types.ValueMetadata {
	keys := map[types.Key]string{constants.Ilk: ilk}
	return types.GetValueMetadata(IlkRate, keys, types.Uint256)
}

func getIlkSpotKey(ilk string) common.Hash {
	return vdbStorage.GetIncrementedKey(getIlkArtKey(ilk), 2)
}

func getIlkSpotMetadata(ilk string) types.ValueMetadata {
	keys := map[types.Key]string{constants.Ilk: ilk}
	return types.GetValueMetadata(IlkSpot, keys, types.Uint256)
}

func getIlkLineKey(ilk string) common.Hash {
	return vdbStorage.GetIncrementedKey(getIlkArtKey(ilk), 3)
}

func getIlkLineMetadata(ilk string) types.ValueMetadata {
	keys := map[types.Key]string{constants.Ilk: ilk}
	return types.GetValueMetadata(IlkLine, keys, types.Uint256)
}

func getIlkDustKey(ilk string) common.Hash {
	return vdbStorage.GetIncrementedKey(getIlkArtKey(ilk), 4)
}

func getIlkDustMetadata(ilk string) types.ValueMetadata {
	keys := map[types.Key]string{constants.Ilk: ilk}
	return types.GetValueMetadata(IlkDust, keys, types.Uint256)
}

func getUrnInkKey(ilk, guy string) common.Hash {
	return vdbStorage.GetKeyForNestedMapping(UrnsMappingIndex, ilk, guy)
}

func getUrnInkMetadata(ilk, guy string) types.ValueMetadata {
	keys := map[types.Key]string{constants.Ilk: ilk, constants.Guy: guy}
	return types.GetValueMetadata(UrnInk, keys, types.Uint256)
}

func getUrnArtKey(ilk, guy string) common.Hash {
	return vdbStorage.GetIncrementedKey(getUrnInkKey(ilk, guy), 1)
}

func getUrnArtMetadata(ilk, guy string) types.ValueMetadata {
	keys := map[types.Key]string{constants.Ilk: ilk, constants.Guy: guy}
	return types.GetValueMetadata(UrnArt, keys, types.Uint256)
}

func getGemKey(ilk, guy string) common.Hash {
	return vdbStorage.GetKeyForNestedMapping(GemsMappingIndex, ilk, guy)
}

func getGemMetadata(ilk, guy string) types.ValueMetadata {
	keys := map[types.Key]string{constants.Ilk: ilk, constants.Guy: guy}
	return types.GetValueMetadata(Gem, keys, types.Uint256)
}

func getDaiKey(guy string) common.Hash {
	return vdbStorage.GetKeyForMapping(DaiMappingIndex, guy)
}

func getDaiMetadata(guy string) types.ValueMetadata {
	keys := map[types.Key]string{constants.Guy: guy}
	return types.GetValueMetadata(Dai, keys, types.Uint256)
}

func getSinKey(guy string) common.Hash {
	return vdbStorage.GetKeyForMapping(SinMappingIndex, guy)
}

func getSinMetadata(guy string) types.ValueMetadata {
	keys := map[types.Key]string{constants.Guy: guy}
	return types.GetValueMetadata(Sin, keys, types.Uint256)
}
