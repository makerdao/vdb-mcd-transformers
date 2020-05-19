package shared

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/backfill/repository"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/vat"
)

func GetUrnArtKey(u repository.Urn) (common.Hash, error) {
	paddedGuy, padErr := utilities.PadAddress(u.Urn)
	if padErr != nil {
		return common.Hash{}, fmt.Errorf("error padding urn identifier: %w", padErr)
	}
	return vat.GetUrnArtKey(u.Ilk, paddedGuy), nil
}

func GetUrnInkKey(u repository.Urn) (common.Hash, error) {
	paddedGuy, padErr := utilities.PadAddress(u.Urn)
	if padErr != nil {
		return common.Hash{}, fmt.Errorf("error padding urn identifier: %w", padErr)
	}
	return vat.GetUrnInkKey(u.Ilk, paddedGuy), nil
}
