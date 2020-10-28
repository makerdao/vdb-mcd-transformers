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
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	vdbStorage "github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
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
	VatKey      = common.HexToHash(vdbStorage.IndexZero)
	VatMetadata = types.ValueMetadata{
		Name: Vat,
		Keys: nil,
		Type: types.Address,
	}

	CdpiKey      = common.HexToHash(vdbStorage.IndexOne)
	CdpiMetadata = types.ValueMetadata{
		Name: Cdpi,
		Keys: nil,
		Type: types.Uint256,
	}

	UrnsMappingIndex  = vdbStorage.IndexTwo
	ListMappingIndex  = vdbStorage.IndexThree
	OwnsMappingIndex  = vdbStorage.IndexFour
	IlksMappingIndex  = vdbStorage.IndexFive
	FirstMappingIndex = vdbStorage.IndexSix
	LastMappingIndex  = vdbStorage.IndexSeven
	CountMappingIndex = vdbStorage.IndexEight
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
	mappings := getStaticMappings()
	mappings, cdpiErr := loader.loadCdpiKeyMappings(mappings)
	if cdpiErr != nil {
		return nil, fmt.Errorf("error adding cdpi key mappings to cdp manager keys loader: %w", cdpiErr)
	}
	mappings, ownsErr := loader.loadOwnsKeyMappings(mappings)
	if ownsErr != nil {
		return nil, fmt.Errorf("error adding owns keys to cdp manager keys loader: %w", ownsErr)
	}
	return mappings, nil
}

func (loader *keysLoader) loadCdpiKeyMappings(mappings map[common.Hash]types.ValueMetadata) (map[common.Hash]types.ValueMetadata, error) {
	cdpis, cdpiErr := loader.storageRepository.GetCdpis()
	if cdpiErr != nil {
		return nil, fmt.Errorf("error getting cdpis: %w", cdpiErr)
	}
	for _, cdpi := range cdpis {
		hexCdpi, hexErr := shared.ConvertIntStringToHex(cdpi)
		if hexErr != nil {
			return nil, fmt.Errorf("error converting int string to hex: %w", hexErr)
		}
		mappings[getUrnsKey(hexCdpi)] = getUrnsMetadata(cdpi)
		mappings[getListPrevKey(hexCdpi)] = getListPrevMetadata(cdpi)
		mappings[getListNextKey(hexCdpi)] = getListNextMetadata(cdpi)
		mappings[getOwnsKey(hexCdpi)] = getOwnsMetadata(cdpi)
		mappings[getIlksKey(hexCdpi)] = getIlksMetadata(cdpi)
	}
	return mappings, nil
}

func (loader *keysLoader) loadOwnsKeyMappings(mappings map[common.Hash]types.ValueMetadata) (map[common.Hash]types.ValueMetadata, error) {
	owners, ownersErr := loader.storageRepository.GetOwners()
	if ownersErr != nil {
		return nil, fmt.Errorf("error getting owners: %w", ownersErr)
	}
	for _, owner := range owners {
		paddedOwner, padErr := utilities.PadAddress(owner)
		if padErr != nil {
			return nil, fmt.Errorf("error padding address: %w", padErr)
		}
		mappings[getFirstKey(paddedOwner)] = getFirstMetadata(owner)
		mappings[getLastKey(paddedOwner)] = getLastMetadata(owner)
		mappings[getCountKey(paddedOwner)] = getCountMetadata(owner)
	}
	return mappings, nil
}

func getStaticMappings() map[common.Hash]types.ValueMetadata {
	mappings := make(map[common.Hash]types.ValueMetadata)
	mappings[VatKey] = VatMetadata
	mappings[CdpiKey] = CdpiMetadata
	return mappings
}

func getUrnsKey(hexCdpi string) common.Hash {
	return vdbStorage.GetKeyForMapping(UrnsMappingIndex, hexCdpi)
}

func getUrnsMetadata(cdpi string) types.ValueMetadata {
	keys := map[types.Key]string{constants.Cdpi: cdpi}
	return types.GetValueMetadata(Urns, keys, types.Address)
}

func getListPrevKey(hexCdpi string) common.Hash {
	return vdbStorage.GetKeyForMapping(ListMappingIndex, hexCdpi)
}

func getListPrevMetadata(cdpi string) types.ValueMetadata {
	keys := map[types.Key]string{constants.Cdpi: cdpi}
	return types.GetValueMetadata(ListPrev, keys, types.Uint256)
}

func getListNextKey(hexCdpi string) common.Hash {
	return vdbStorage.GetIncrementedKey(getListPrevKey(hexCdpi), 1)
}

func getListNextMetadata(cdpi string) types.ValueMetadata {
	keys := map[types.Key]string{constants.Cdpi: cdpi}
	return types.GetValueMetadata(ListNext, keys, types.Uint256)
}

func getOwnsKey(hexCdpi string) common.Hash {
	return vdbStorage.GetKeyForMapping(OwnsMappingIndex, hexCdpi)
}

func getOwnsMetadata(cdpi string) types.ValueMetadata {
	keys := map[types.Key]string{constants.Cdpi: cdpi}
	return types.GetValueMetadata(Owns, keys, types.Address)
}

func getIlksKey(hexCdpi string) common.Hash {
	return vdbStorage.GetKeyForMapping(IlksMappingIndex, hexCdpi)
}

func getIlksMetadata(cdpi string) types.ValueMetadata {
	keys := map[types.Key]string{constants.Cdpi: cdpi}
	return types.GetValueMetadata(Ilks, keys, types.Bytes32)
}

func getFirstKey(ownerAddress string) common.Hash {
	return vdbStorage.GetKeyForMapping(FirstMappingIndex, ownerAddress)
}

func getFirstMetadata(ownerAddress string) types.ValueMetadata {
	keys := map[types.Key]string{constants.Owner: ownerAddress}
	return types.GetValueMetadata(First, keys, types.Uint256)
}

func getLastKey(ownerAddress string) common.Hash {
	return vdbStorage.GetKeyForMapping(LastMappingIndex, ownerAddress)
}

func getLastMetadata(ownerAddress string) types.ValueMetadata {
	keys := map[types.Key]string{constants.Owner: ownerAddress}
	return types.GetValueMetadata(Last, keys, types.Uint256)
}

func getCountKey(ownerAddress string) common.Hash {
	return vdbStorage.GetKeyForMapping(CountMappingIndex, ownerAddress)
}

func getCountMetadata(ownerAddress string) types.ValueMetadata {
	keys := map[types.Key]string{constants.Owner: ownerAddress}
	return types.GetValueMetadata(Count, keys, types.Uint256)
}
