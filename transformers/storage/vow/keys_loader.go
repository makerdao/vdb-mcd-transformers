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
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities/wards"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	vdbStorage "github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

const (
	Vat        = "vat"
	Flapper    = "flapper"
	Flopper    = "flopper"
	SinMapping = "sin"
	SinInteger = "Sin"
	Ash        = "Ash"
	Wait       = "wait"
	Dump       = "dump"
	Sump       = "sump"
	Bump       = "bump"
	Hump       = "hump"
	Live       = "live"
)

var (
	VatKey      = common.HexToHash(vdbStorage.IndexOne)
	VatMetadata = types.ValueMetadata{
		Name: Vat,
		Keys: nil,
		Type: types.Address,
	}

	FlapperKey      = common.HexToHash(vdbStorage.IndexTwo)
	FlapperMetadata = types.ValueMetadata{
		Name: Flapper,
		Keys: nil,
		Type: types.Address,
	}

	FlopperKey      = common.HexToHash(vdbStorage.IndexThree)
	FlopperMetadata = types.ValueMetadata{
		Name: Flopper,
		Keys: nil,
		Type: types.Address,
	}

	SinMappingIndex = vdbStorage.IndexFour

	SinIntegerKey      = common.HexToHash(vdbStorage.IndexFive)
	SinIntegerMetadata = types.ValueMetadata{
		Name: SinInteger,
		Keys: nil,
		Type: types.Uint256,
	}

	AshKey      = common.HexToHash(vdbStorage.IndexSix)
	AshMetadata = types.ValueMetadata{
		Name: Ash,
		Keys: nil,
		Type: types.Uint256,
	}

	WaitKey      = common.HexToHash(vdbStorage.IndexSeven)
	WaitMetadata = types.ValueMetadata{
		Name: Wait,
		Keys: nil,
		Type: types.Uint256,
	}

	DumpKey      = common.HexToHash(vdbStorage.IndexEight)
	DumpMetadata = types.ValueMetadata{
		Name: Dump,
		Keys: nil,
		Type: types.Uint256,
	}

	SumpKey      = common.HexToHash(vdbStorage.IndexNine)
	SumpMetadata = types.ValueMetadata{
		Name: Sump,
		Keys: nil,
		Type: types.Uint256,
	}

	BumpKey      = common.HexToHash(vdbStorage.IndexTen)
	BumpMetadata = types.ValueMetadata{
		Name: Bump,
		Keys: nil,
		Type: types.Uint256,
	}

	HumpKey      = common.HexToHash(vdbStorage.IndexEleven)
	HumpMetadata = types.ValueMetadata{
		Name: Hump,
		Keys: nil,
		Type: types.Uint256,
	}

	LiveKey      = common.HexToHash(vdbStorage.IndexTwelve)
	LiveMetadata = types.ValueMetadata{
		Name: Live,
		Keys: nil,
		Type: types.Uint256,
	}
)

type keysLoader struct {
	storageRepository mcdStorage.IMakerStorageRepository
	contractAddress   string
}

func NewKeysLoader(storageRepository mcdStorage.IMakerStorageRepository, contractAddress string) storage.KeysLoader {
	return &keysLoader{storageRepository: storageRepository, contractAddress: contractAddress}
}

func (loader *keysLoader) LoadMappings() (map[common.Hash]types.ValueMetadata, error) {
	mappings := getStaticMappings()
	return loader.addDynamicMappings(mappings)
}

func (loader *keysLoader) SetDB(db *postgres.DB) {
	loader.storageRepository.SetDB(db)
}

func (loader *keysLoader) addDynamicMappings(mappings map[common.Hash]types.ValueMetadata) (map[common.Hash]types.ValueMetadata, error) {
	mappings, wardsErr := loader.addWardsKeys(mappings)
	if wardsErr != nil {
		return nil, fmt.Errorf("error adding wards keys to vow keys loader: %w", wardsErr)
	}
	mappings, sinErr := loader.addVowSinKeys(mappings)
	if sinErr != nil {
		return nil, fmt.Errorf("error adding sin keys to vow keys loader: %w", sinErr)
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

func (loader *keysLoader) addVowSinKeys(mappings map[common.Hash]types.ValueMetadata) (map[common.Hash]types.ValueMetadata, error) {
	sinKeys, getErr := loader.storageRepository.GetVowSinKeys()
	if getErr != nil {
		return nil, fmt.Errorf("error getting vow sin keys: %w", getErr)
	}
	for _, timestamp := range sinKeys {
		hexTimestamp, convertErr := shared.ConvertIntStringToHex(timestamp)
		if convertErr != nil {
			return nil, fmt.Errorf("error converting int string to hex: %w", convertErr)
		}
		mappings[getSinKey(hexTimestamp)] = getSinMetadata(timestamp)
	}
	return mappings, nil
}

func getStaticMappings() map[common.Hash]types.ValueMetadata {
	mappings := make(map[common.Hash]types.ValueMetadata)
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
	mappings[LiveKey] = LiveMetadata
	return mappings
}

func getSinKey(hexTimestamp string) common.Hash {
	return vdbStorage.GetKeyForMapping(SinMappingIndex, hexTimestamp)
}

func getSinMetadata(timestamp string) types.ValueMetadata {
	keys := map[types.Key]string{constants.Timestamp: timestamp}
	return types.GetValueMetadata(SinMapping, keys, types.Uint256)
}
