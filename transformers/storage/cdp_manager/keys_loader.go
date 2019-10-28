// VulcanizeDB
// Copyright © 2019 Vulcanize

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

package cdp_manager

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	mcdStorage "github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/mcd_transformers/transformers/storage/utilities"
	"github.com/vulcanize/vulcanizedb/libraries/shared/factories/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

const (
	Vat      = "vat"
	Cdpi     = "cdpi"
	Urns     = "urns"
	ListPrev = "prev"
	ListNext = "next"
	Owns     = "owns"
	Ilks     = "ilks"
	First    = "first"
	Last     = "last"
	Count    = "count"
)

var (
	VatKey      = common.HexToHash(utils.IndexZero)
	VatMetadata = utils.StorageValueMetadata{
		Name: Vat,
		Keys: nil,
		Type: utils.Address,
	}

	CdpiKey      = common.HexToHash(utils.IndexOne)
	CdpiMetadata = utils.StorageValueMetadata{
		Name: Cdpi,
		Keys: nil,
		Type: utils.Uint256,
	}

	UrnsMappingIndex  = utils.IndexTwo
	ListMappingIndex  = utils.IndexThree
	OwnsMappingIndex  = utils.IndexFour
	IlksMappingIndex  = utils.IndexFive
	FirstMappingIndex = utils.IndexSix
	LastMappingIndex  = utils.IndexSeven
	CountMappingIndex = utils.IndexEight
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

func (loader *keysLoader) LoadMappings() (map[common.Hash]utils.StorageValueMetadata, error) {
	mappings := loadStaticMappings()
	mappings, cdpiErr := loader.loadCdpiKeyMappings(mappings)
	if cdpiErr != nil {
		return nil, cdpiErr
	}
	return loader.loadOwnsKeyMappings(mappings)
}

func (loader *keysLoader) loadCdpiKeyMappings(mappings map[common.Hash]utils.StorageValueMetadata) (map[common.Hash]utils.StorageValueMetadata, error) {
	cdpis, cdpiErr := loader.storageRepository.GetCdpis()
	if cdpiErr != nil {
		return nil, cdpiErr
	}
	for _, cdpi := range cdpis {
		hexCdpi, hexErr := shared.ConvertIntStringToHex(cdpi)
		if hexErr != nil {
			return nil, hexErr
		}
		mappings[getUrnsKey(hexCdpi)] = getUrnsMetadata(cdpi)
		mappings[getListPrevKey(hexCdpi)] = getListPrevMetadata(cdpi)
		mappings[getListNextKey(hexCdpi)] = getListNextMetadata(cdpi)
		mappings[getOwnsKey(hexCdpi)] = getOwnsMetadata(cdpi)
		mappings[getIlksKey(hexCdpi)] = getIlksMetadata(cdpi)
	}
	return mappings, nil
}

func (loader *keysLoader) loadOwnsKeyMappings(mappings map[common.Hash]utils.StorageValueMetadata) (map[common.Hash]utils.StorageValueMetadata, error) {
	owners, ownersErr := loader.storageRepository.GetOwners()
	if ownersErr != nil {
		return nil, ownersErr
	}
	for _, owner := range owners {
		paddedOwner, padErr := utilities.PadAddress(owner)
		if padErr != nil {
			return nil, padErr
		}
		mappings[getFirstKey(paddedOwner)] = getFirstMetadata(owner)
		mappings[getLastKey(paddedOwner)] = getLastMetadata(owner)
		mappings[getCountKey(paddedOwner)] = getCountMetadata(owner)
	}
	return mappings, nil
}

func loadStaticMappings() map[common.Hash]utils.StorageValueMetadata {
	mappings := make(map[common.Hash]utils.StorageValueMetadata)
	mappings[VatKey] = VatMetadata
	mappings[CdpiKey] = CdpiMetadata
	return mappings
}

func getUrnsKey(hexCdpi string) common.Hash {
	return utils.GetStorageKeyForMapping(UrnsMappingIndex, hexCdpi)
}

func getUrnsMetadata(cdpi string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Cdpi: cdpi}
	return utils.GetStorageValueMetadata(Urns, keys, utils.Address)
}

func getListPrevKey(hexCdpi string) common.Hash {
	return utils.GetStorageKeyForMapping(ListMappingIndex, hexCdpi)
}

func getListPrevMetadata(cdpi string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Cdpi: cdpi}
	return utils.GetStorageValueMetadata(ListPrev, keys, utils.Uint256)
}

func getListNextKey(hexCdpi string) common.Hash {
	return utils.GetIncrementedStorageKey(getListPrevKey(hexCdpi), 1)
}

func getListNextMetadata(cdpi string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Cdpi: cdpi}
	return utils.GetStorageValueMetadata(ListNext, keys, utils.Uint256)
}

func getOwnsKey(hexCdpi string) common.Hash {
	return utils.GetStorageKeyForMapping(OwnsMappingIndex, hexCdpi)
}

func getOwnsMetadata(cdpi string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Cdpi: cdpi}
	return utils.GetStorageValueMetadata(Owns, keys, utils.Address)
}

func getIlksKey(hexCdpi string) common.Hash {
	return utils.GetStorageKeyForMapping(IlksMappingIndex, hexCdpi)
}

func getIlksMetadata(cdpi string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Cdpi: cdpi}
	return utils.GetStorageValueMetadata(Ilks, keys, utils.Bytes32)
}

func getFirstKey(ownerAddress string) common.Hash {
	return utils.GetStorageKeyForMapping(FirstMappingIndex, ownerAddress)
}

func getFirstMetadata(ownerAddress string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Owner: ownerAddress}
	return utils.GetStorageValueMetadata(First, keys, utils.Uint256)
}

func getLastKey(ownerAddress string) common.Hash {
	return utils.GetStorageKeyForMapping(LastMappingIndex, ownerAddress)
}

func getLastMetadata(ownerAddress string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Owner: ownerAddress}
	return utils.GetStorageValueMetadata(Last, keys, utils.Uint256)
}

func getCountKey(ownerAddress string) common.Hash {
	return utils.GetStorageKeyForMapping(CountMappingIndex, ownerAddress)
}

func getCountMetadata(ownerAddress string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Owner: ownerAddress}
	return utils.GetStorageValueMetadata(Count, keys, utils.Uint256)
}
