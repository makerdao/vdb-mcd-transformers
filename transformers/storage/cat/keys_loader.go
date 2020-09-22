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
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities/wards"
	vdbStorage "github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
)

var (
	IlksMappingIndex = vdbStorage.IndexOne // bytes32 => flip address; chop (ray), lump (wad) uint256

	LiveKey      = common.HexToHash(vdbStorage.IndexTwo)
	LiveMetadata = types.GetValueMetadata(Live, nil, types.Uint256)

	VatKey      = common.HexToHash(vdbStorage.IndexThree)
	VatMetadata = types.GetValueMetadata(Vat, nil, types.Address)

	VowKey      = common.HexToHash(vdbStorage.IndexFour)
	VowMetadata = types.GetValueMetadata(Vow, nil, types.Address)
)

func LoadSharedMappings(mappings map[common.Hash]types.ValueMetadata, address string, repository mcdStorage.IMakerStorageRepository) (map[common.Hash]types.ValueMetadata, error) {
	mappings = loadSharedStaticMappings(mappings)
	mappings, wardsErr := addWardsKeys(mappings, address, repository)
	if wardsErr != nil {
		return nil, fmt.Errorf("error adding wards keys to cat keys loader: %w", wardsErr)
	}
	return mappings, nil
}

func addWardsKeys(mappings map[common.Hash]types.ValueMetadata, address string, repository mcdStorage.IMakerStorageRepository) (map[common.Hash]types.ValueMetadata, error) {
	addresses, err := repository.GetWardsAddresses(address)
	if err != nil {
		return nil, fmt.Errorf("error getting wards addresses: %w", err)
	}
	return wards.AddWardsKeys(mappings, addresses)
}

func loadSharedStaticMappings(mappings map[common.Hash]types.ValueMetadata) map[common.Hash]types.ValueMetadata {
	mappings[LiveKey] = LiveMetadata
	mappings[VatKey] = VatMetadata
	mappings[VowKey] = VowMetadata
	return mappings
}

func GetIlkFlipKey(ilk string) common.Hash {
	return vdbStorage.GetKeyForMapping(IlksMappingIndex, ilk)
}

func GetIlkFlipMetadata(ilk string) types.ValueMetadata {
	keys := map[types.Key]string{constants.Ilk: ilk}
	return types.GetValueMetadata(IlkFlip, keys, types.Address)
}

func GetIlkChopKey(ilk string) common.Hash {
	return vdbStorage.GetIncrementedKey(GetIlkFlipKey(ilk), 1)
}

func GetIlkChopMetadata(ilk string) types.ValueMetadata {
	keys := map[types.Key]string{constants.Ilk: ilk}
	return types.GetValueMetadata(IlkChop, keys, types.Uint256)
}
