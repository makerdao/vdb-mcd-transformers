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
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	s2 "github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/mcd_transformers/transformers/storage/utilities"
)

const (
	CdpManagerVat      = "vat"
	CdpManagerCdpi     = "cdpi"
	CdpManagerUrns     = "urns"
	CdpManagerListPrev = "prev"
	CdpManagerListNext = "next"
	CdpManagerOwns     = "owns"
	CdpManagerIlks     = "ilks"
	CdpManagerFirst    = "first"
	CdpManagerLast     = "last"
	CdpManagerCount    = "count"
)

var (
	VatKey      = common.HexToHash(storage.IndexZero)
	VatMetadata = utils.StorageValueMetadata{
		Name: CdpManagerVat,
		Keys: nil,
		Type: utils.Address,
	}

	CdpiKey      = common.HexToHash(storage.IndexOne)
	CdpiMetadata = utils.StorageValueMetadata{
		Name: CdpManagerCdpi,
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

type CdpManagerMappings struct {
	StorageRepository s2.IMakerStorageRepository
	mappings          map[common.Hash]utils.StorageValueMetadata
}

func (mappings *CdpManagerMappings) SetDB(db *postgres.DB) {
	mappings.StorageRepository.SetDB(db)
}

func (mappings *CdpManagerMappings) Lookup(key common.Hash) (utils.StorageValueMetadata, error) {
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

func (mappings *CdpManagerMappings) loadMappings() error {
	mappings.mappings = loadStaticMappings()
	cdpiErr := mappings.loadCdpiKeyMappings()
	if cdpiErr != nil {
		return cdpiErr
	}
	ownsErr := mappings.loadOwnsKeyMappings()
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

func (mappings *CdpManagerMappings) loadCdpiKeyMappings() error {
	cdpis, cdpiErr := mappings.StorageRepository.GetCdpis()
	if cdpiErr != nil {
		return cdpiErr
	}
	for _, cdpi := range cdpis {
		hexCdpi, hexErr := shared.ConvertIntStringToHex(cdpi)
		if hexErr != nil {
			return hexErr
		}
		mappings.mappings[getUrnsKey(hexCdpi)] = getUrnsMetadata(cdpi)
		mappings.mappings[getListPrevKey(hexCdpi)] = getListPrevMetadata(cdpi)
		mappings.mappings[getListNextKey(hexCdpi)] = getListNextMetadata(cdpi)
		mappings.mappings[getOwnsKey(hexCdpi)] = getOwnsMetadata(cdpi)
		mappings.mappings[getIlksKey(hexCdpi)] = getIlksMetadata(cdpi)
	}
	return nil
}

func (mappings *CdpManagerMappings) loadOwnsKeyMappings() error {
	owners, ownersErr := mappings.StorageRepository.GetOwners()
	if ownersErr != nil {
		return ownersErr
	}
	for _, owner := range owners {
		paddedOwner, padErr := utilities.PadAddress(owner)
		if padErr != nil {
			return padErr
		}
		mappings.mappings[getFirstKey(paddedOwner)] = getFirstMetadata(owner)
		mappings.mappings[getLastKey(paddedOwner)] = getLastMetadata(owner)
		mappings.mappings[getCountKey(paddedOwner)] = getCountMetadata(owner)
	}
	return nil
}

func getUrnsKey(hexCdpi string) common.Hash {
	return storage.GetMapping(UrnsMappingIndex, hexCdpi)
}

func getUrnsMetadata(cdpi string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Cdpi: cdpi}
	return utils.GetStorageValueMetadata(CdpManagerUrns, keys, utils.Address)
}

func getListPrevKey(hexCdpi string) common.Hash {
	return storage.GetMapping(ListMappingIndex, hexCdpi)
}

func getListPrevMetadata(cdpi string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Cdpi: cdpi}
	return utils.GetStorageValueMetadata(CdpManagerListPrev, keys, utils.Uint256)
}

func getListNextKey(hexCdpi string) common.Hash {
	return storage.GetIncrementedKey(getListPrevKey(hexCdpi), 1)
}

func getListNextMetadata(cdpi string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Cdpi: cdpi}
	return utils.GetStorageValueMetadata(CdpManagerListNext, keys, utils.Uint256)
}

func getOwnsKey(hexCdpi string) common.Hash {
	return storage.GetMapping(OwnsMappingIndex, hexCdpi)
}

func getOwnsMetadata(cdpi string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Cdpi: cdpi}
	return utils.GetStorageValueMetadata(CdpManagerOwns, keys, utils.Address)
}

func getIlksKey(hexCdpi string) common.Hash {
	return storage.GetMapping(IlksMappingIndex, hexCdpi)
}

func getIlksMetadata(cdpi string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Cdpi: cdpi}
	return utils.GetStorageValueMetadata(CdpManagerIlks, keys, utils.Bytes32)
}

func getFirstKey(ownerAddress string) common.Hash {
	return storage.GetMapping(FirstMappingIndex, ownerAddress)
}

func getFirstMetadata(ownerAddress string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Owner: ownerAddress}
	return utils.GetStorageValueMetadata(CdpManagerFirst, keys, utils.Uint256)
}

func getLastKey(ownerAddress string) common.Hash {
	return storage.GetMapping(LastMappingIndex, ownerAddress)
}

func getLastMetadata(ownerAddress string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Owner: ownerAddress}
	return utils.GetStorageValueMetadata(CdpManagerLast, keys, utils.Uint256)
}

func getCountKey(ownerAddress string) common.Hash {
	return storage.GetMapping(CountMappingIndex, ownerAddress)
}

func getCountMetadata(ownerAddress string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Owner: ownerAddress}
	return utils.GetStorageValueMetadata(CdpManagerCount, keys, utils.Uint256)
}
