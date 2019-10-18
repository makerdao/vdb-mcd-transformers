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
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage"
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
	VatKey      = common.HexToHash(storage.IndexZero)
	VatMetadata = utils.StorageValueMetadata{
		Name: Vat,
		Keys: nil,
		Type: utils.Address,
	}

	CdpiKey      = common.HexToHash(storage.IndexOne)
	CdpiMetadata = utils.StorageValueMetadata{
		Name: Cdpi,
		Keys: nil,
		Type: utils.Uint256,
	}

	UrnsMappingIndex  = storage.IndexTwo
	ListMappingIndex  = storage.IndexThree
	OwnsMappingIndex  = storage.IndexFour
	IlksMappingIndex  = storage.IndexFive
	FirstMappingIndex = storage.IndexSix
	LastMappingIndex  = storage.IndexSeven
	CountMappingIndex = storage.IndexEight
)

type StorageKeysLookup struct {
	StorageRepository mcdStorage.IMakerStorageRepository
	mappings          map[common.Hash]utils.StorageValueMetadata
}

func (lookup *StorageKeysLookup) SetDB(db *postgres.DB) {
	lookup.StorageRepository.SetDB(db)
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
	cdpiErr := lookup.loadCdpiKeyMappings()
	if cdpiErr != nil {
		return cdpiErr
	}
	ownsErr := lookup.loadOwnsKeyMappings()
	if ownsErr != nil {
		return ownsErr
	}
	return nil
}

func loadStaticMappings() map[common.Hash]utils.StorageValueMetadata {
	mappings := make(map[common.Hash]utils.StorageValueMetadata)
	mappings[VatKey] = VatMetadata
	mappings[CdpiKey] = CdpiMetadata
	return mappings
}

func (lookup *StorageKeysLookup) loadCdpiKeyMappings() error {
	cdpis, cdpiErr := lookup.StorageRepository.GetCdpis()
	if cdpiErr != nil {
		return cdpiErr
	}
	for _, cdpi := range cdpis {
		hexCdpi, hexErr := shared.ConvertIntStringToHex(cdpi)
		if hexErr != nil {
			return hexErr
		}
		lookup.mappings[getUrnsKey(hexCdpi)] = getUrnsMetadata(cdpi)
		lookup.mappings[getListPrevKey(hexCdpi)] = getListPrevMetadata(cdpi)
		lookup.mappings[getListNextKey(hexCdpi)] = getListNextMetadata(cdpi)
		lookup.mappings[getOwnsKey(hexCdpi)] = getOwnsMetadata(cdpi)
		lookup.mappings[getIlksKey(hexCdpi)] = getIlksMetadata(cdpi)
	}
	return nil
}

func (lookup *StorageKeysLookup) loadOwnsKeyMappings() error {
	owners, ownersErr := lookup.StorageRepository.GetOwners()
	if ownersErr != nil {
		return ownersErr
	}
	for _, owner := range owners {
		paddedOwner, padErr := utilities.PadAddress(owner)
		if padErr != nil {
			return padErr
		}
		lookup.mappings[getFirstKey(paddedOwner)] = getFirstMetadata(owner)
		lookup.mappings[getLastKey(paddedOwner)] = getLastMetadata(owner)
		lookup.mappings[getCountKey(paddedOwner)] = getCountMetadata(owner)
	}
	return nil
}

func getUrnsKey(hexCdpi string) common.Hash {
	return storage.GetMapping(UrnsMappingIndex, hexCdpi)
}

func getUrnsMetadata(cdpi string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Cdpi: cdpi}
	return utils.GetStorageValueMetadata(Urns, keys, utils.Address)
}

func getListPrevKey(hexCdpi string) common.Hash {
	return storage.GetMapping(ListMappingIndex, hexCdpi)
}

func getListPrevMetadata(cdpi string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Cdpi: cdpi}
	return utils.GetStorageValueMetadata(ListPrev, keys, utils.Uint256)
}

func getListNextKey(hexCdpi string) common.Hash {
	return storage.GetIncrementedKey(getListPrevKey(hexCdpi), 1)
}

func getListNextMetadata(cdpi string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Cdpi: cdpi}
	return utils.GetStorageValueMetadata(ListNext, keys, utils.Uint256)
}

func getOwnsKey(hexCdpi string) common.Hash {
	return storage.GetMapping(OwnsMappingIndex, hexCdpi)
}

func getOwnsMetadata(cdpi string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Cdpi: cdpi}
	return utils.GetStorageValueMetadata(Owns, keys, utils.Address)
}

func getIlksKey(hexCdpi string) common.Hash {
	return storage.GetMapping(IlksMappingIndex, hexCdpi)
}

func getIlksMetadata(cdpi string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Cdpi: cdpi}
	return utils.GetStorageValueMetadata(Ilks, keys, utils.Bytes32)
}

func getFirstKey(ownerAddress string) common.Hash {
	return storage.GetMapping(FirstMappingIndex, ownerAddress)
}

func getFirstMetadata(ownerAddress string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Owner: ownerAddress}
	return utils.GetStorageValueMetadata(First, keys, utils.Uint256)
}

func getLastKey(ownerAddress string) common.Hash {
	return storage.GetMapping(LastMappingIndex, ownerAddress)
}

func getLastMetadata(ownerAddress string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Owner: ownerAddress}
	return utils.GetStorageValueMetadata(Last, keys, utils.Uint256)
}

func getCountKey(ownerAddress string) common.Hash {
	return storage.GetMapping(CountMappingIndex, ownerAddress)
}

func getCountMetadata(ownerAddress string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Owner: ownerAddress}
	return utils.GetStorageValueMetadata(Count, keys, utils.Uint256)
}
