package shared

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/backfill/repository"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage"
)

func GetIlkArtKey(ilk string) common.Hash {
	return storage.GetKeyForMapping(storage.IndexTwo, ilk)
}

func GetUrnArtKey(u repository.Urn) (common.Hash, error) {
	inkKey, err := GetUrnInkKey(u)
	if err != nil {
		return common.Hash{}, err
	}
	artKey := storage.GetIncrementedKey(inkKey, 1)
	return artKey, nil
}

func GetUrnInkKey(u repository.Urn) (common.Hash, error) {
	paddedGuy, padErr := utilities.PadAddress(u.Urn)
	if padErr != nil {
		return common.Hash{}, fmt.Errorf("error padding urn identifier: %w", padErr)
	}
	inkKey := storage.GetKeyForNestedMapping(storage.IndexThree, u.Ilk, paddedGuy)
	return inkKey, nil
}
