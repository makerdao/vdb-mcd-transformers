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

package cat

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	vdbStorage "github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

const (
	Live = "live"
	Vat  = "vat"
	Vow  = "vow"

	Wards   = "wards"
	IlkFlip = "flip"
	IlkChop = "chop"
	IlkLump = "lump"
)

var (
	WardsMappingIndex = vdbStorage.IndexZero

	IlksMappingIndex = vdbStorage.IndexOne // bytes32 => flip address; chop (ray), lump (wad) uint256

	LiveKey      = common.HexToHash(vdbStorage.IndexTwo)
	LiveMetadata = vdbStorage.GetValueMetadata(Live, nil, vdbStorage.Uint256)

	VatKey      = common.HexToHash(vdbStorage.IndexThree)
	VatMetadata = vdbStorage.GetValueMetadata(Vat, nil, vdbStorage.Address)

	VowKey      = common.HexToHash(vdbStorage.IndexFour)
	VowMetadata = vdbStorage.GetValueMetadata(Vow, nil, vdbStorage.Address)
)

type keysLoader struct {
	storageRepository mcdStorage.IMakerStorageRepository
	contractAddress   string
}

func NewKeysLoader(storageRepository mcdStorage.IMakerStorageRepository, contractAddress string) storage.KeysLoader {
	return &keysLoader{storageRepository: storageRepository, contractAddress: contractAddress}
}

func (loader *keysLoader) SetDB(db *postgres.DB) {
	loader.storageRepository.SetDB(db)
}

func (loader *keysLoader) LoadMappings() (map[common.Hash]vdbStorage.ValueMetadata, error) {
	mappings := loadStaticMappings()
	mappings, ilkErr := loader.addIlkKeys(mappings)
	if ilkErr != nil {
		return nil, ilkErr
	}
	mappings, wardsErr := loader.addWardsKeys(mappings)
	if wardsErr != nil {
		return nil, wardsErr
	}
	return mappings, nil
}

func (loader *keysLoader) addIlkKeys(mappings map[common.Hash]vdbStorage.ValueMetadata) (map[common.Hash]vdbStorage.ValueMetadata, error) {
	ilks, err := loader.storageRepository.GetIlks()
	if err != nil {
		return nil, err
	}
	for _, ilk := range ilks {
		mappings[getIlkFlipKey(ilk)] = getIlkFlipMetadata(ilk)
		mappings[getIlkChopKey(ilk)] = getIlkChopMetadata(ilk)
		mappings[getIlkLumpKey(ilk)] = getIlkLumpMetadata(ilk)
	}
	return mappings, nil
}

func (loader *keysLoader) addWardsKeys(mappings map[common.Hash]vdbStorage.ValueMetadata) (map[common.Hash]vdbStorage.ValueMetadata, error) {
	addresses, err := loader.storageRepository.GetWardsAddresses(loader.contractAddress)
	if err != nil {
		return nil, err
	}
	for _, address := range addresses {
		paddedAddress, padErr := utilities.PadAddress(address)
		if padErr != nil {
			return nil, padErr
		}
		mappings[getWardsKey(paddedAddress)] = getWardsMetadata(address)
	}
	return mappings, nil
}

func loadStaticMappings() map[common.Hash]vdbStorage.ValueMetadata {
	mappings := make(map[common.Hash]vdbStorage.ValueMetadata)
	mappings[LiveKey] = LiveMetadata
	mappings[VatKey] = VatMetadata
	mappings[VowKey] = VowMetadata
	return mappings
}

func getIlkFlipKey(ilk string) common.Hash {
	return vdbStorage.GetKeyForMapping(IlksMappingIndex, ilk)
}

func getIlkFlipMetadata(ilk string) vdbStorage.ValueMetadata {
	keys := map[vdbStorage.Key]string{constants.Ilk: ilk}
	return vdbStorage.GetValueMetadata(IlkFlip, keys, vdbStorage.Address)
}

func getIlkChopKey(ilk string) common.Hash {
	return vdbStorage.GetIncrementedKey(getIlkFlipKey(ilk), 1)
}

func getIlkChopMetadata(ilk string) vdbStorage.ValueMetadata {
	keys := map[vdbStorage.Key]string{constants.Ilk: ilk}
	return vdbStorage.GetValueMetadata(IlkChop, keys, vdbStorage.Uint256)
}

func getIlkLumpKey(ilk string) common.Hash {
	return vdbStorage.GetIncrementedKey(getIlkFlipKey(ilk), 2)
}

func getIlkLumpMetadata(ilk string) vdbStorage.ValueMetadata {
	keys := map[vdbStorage.Key]string{constants.Ilk: ilk}
	return vdbStorage.GetValueMetadata(IlkLump, keys, vdbStorage.Uint256)
}

func getWardsKey(address string) common.Hash {
	return vdbStorage.GetKeyForMapping(WardsMappingIndex, address)
}

func getWardsMetadata(user string) vdbStorage.ValueMetadata {
	keys := map[vdbStorage.Key]string{constants.User: user}
	return vdbStorage.GetValueMetadata(Wards, keys, vdbStorage.Uint256)
}
