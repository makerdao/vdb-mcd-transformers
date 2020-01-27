package wards

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
)

var (
	Wards             = "wards"
	WardsMappingIndex = storage.IndexZero
)

func AddWardsKeys(mappings map[common.Hash]types.ValueMetadata, addresses []string) (map[common.Hash]types.ValueMetadata, error) {
	for _, address := range addresses {
		paddedAddress, padErr := utilities.PadAddress(address)
		if padErr != nil {
			return nil, padErr
		}
		mappings[getWardsKey(paddedAddress)] = getWardsMetadata(address)
	}
	return mappings, nil
}

func getWardsKey(address string) common.Hash {
	return storage.GetKeyForMapping(WardsMappingIndex, address)
}

func getWardsMetadata(user string) types.ValueMetadata {
	keys := map[types.Key]string{constants.User: user}
	return types.GetValueMetadata(Wards, keys, types.Uint256)
}
