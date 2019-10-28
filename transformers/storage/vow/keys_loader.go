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
	"github.com/vulcanize/vulcanizedb/libraries/shared/factories/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
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
)

var (
	VatKey      = common.HexToHash(utils.IndexOne)
	VatMetadata = utils.StorageValueMetadata{
		Name: Vat,
		Keys: nil,
		Type: utils.Address,
	}

	FlapperKey      = common.HexToHash(utils.IndexTwo)
	FlapperMetadata = utils.StorageValueMetadata{
		Name: Flapper,
		Keys: nil,
		Type: utils.Address,
	}

	FlopperKey      = common.HexToHash(utils.IndexThree)
	FlopperMetadata = utils.StorageValueMetadata{
		Name: Flopper,
		Keys: nil,
		Type: utils.Address,
	}

	SinMappingIndex = utils.IndexFour

	SinIntegerKey      = common.HexToHash(utils.IndexFive)
	SinIntegerMetadata = utils.StorageValueMetadata{
		Name: SinInteger,
		Keys: nil,
		Type: utils.Uint256,
	}

	AshKey      = common.HexToHash(utils.IndexSix)
	AshMetadata = utils.StorageValueMetadata{
		Name: Ash,
		Keys: nil,
		Type: utils.Uint256,
	}

	WaitKey      = common.HexToHash(utils.IndexSeven)
	WaitMetadata = utils.StorageValueMetadata{
		Name: Wait,
		Keys: nil,
		Type: utils.Uint256,
	}

	DumpKey      = common.HexToHash(utils.IndexEight)
	DumpMetadata = utils.StorageValueMetadata{
		Name: Dump,
		Keys: nil,
		Type: utils.Uint256,
	}

	SumpKey      = common.HexToHash(utils.IndexNine)
	SumpMetadata = utils.StorageValueMetadata{
		Name: Sump,
		Keys: nil,
		Type: utils.Uint256,
	}

	BumpKey      = common.HexToHash(utils.IndexTen)
	BumpMetadata = utils.StorageValueMetadata{
		Name: Bump,
		Keys: nil,
		Type: utils.Uint256,
	}

	HumpKey      = common.HexToHash(utils.IndexEleven)
	HumpMetadata = utils.StorageValueMetadata{
		Name: Hump,
		Keys: nil,
		Type: utils.Uint256,
	}
)

type keysLoader struct {
	storageRepository mcdStorage.IMakerStorageRepository
}

func NewKeysLoader(storageRepository mcdStorage.IMakerStorageRepository) storage.KeysLoader {
	return &keysLoader{storageRepository: storageRepository}
}

func (loader *keysLoader) LoadMappings() (map[common.Hash]utils.StorageValueMetadata, error) {
	mappings := addStaticMappings(make(map[common.Hash]utils.StorageValueMetadata))
	return loader.addDynamicMappings(mappings)
}

func (loader *keysLoader) SetDB(db *postgres.DB) {
	loader.storageRepository.SetDB(db)
}

func (loader *keysLoader) addDynamicMappings(mappings map[common.Hash]utils.StorageValueMetadata) (map[common.Hash]utils.StorageValueMetadata, error) {
	sinKeys, getErr := loader.storageRepository.GetVowSinKeys()
	if getErr != nil {
		return nil, getErr
	}
	for _, timestamp := range sinKeys {
		hexTimestamp, convertErr := shared.ConvertIntStringToHex(timestamp)
		if convertErr != nil {
			return nil, convertErr
		}
		mappings[getSinKey(hexTimestamp)] = getSinMetadata(timestamp)
	}
	return mappings, nil
}

func addStaticMappings(mappings map[common.Hash]utils.StorageValueMetadata) map[common.Hash]utils.StorageValueMetadata {
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

func getSinKey(hexTimestamp string) common.Hash {
	return utils.GetStorageKeyForMapping(SinMappingIndex, hexTimestamp)
}

func getSinMetadata(timestamp string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Timestamp: timestamp}
	return utils.GetStorageValueMetadata(SinMapping, keys, utils.Uint256)
}
