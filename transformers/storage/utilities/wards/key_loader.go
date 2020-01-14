package wards

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage"
)

var (
	Wards             = "wards"
	WardsMappingIndex = storage.IndexZero
)

func AddWardsKeys(mappings map[common.Hash]storage.ValueMetadata, addresses []string) (map[common.Hash]storage.ValueMetadata, error) {
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

func getWardsMetadata(user string) storage.ValueMetadata {
	keys := map[storage.Key]string{constants.User: user}
	return storage.GetValueMetadata(Wards, keys, storage.Uint256)
}
