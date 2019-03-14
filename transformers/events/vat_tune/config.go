package vat_tune

import (
	shared_t "github.com/vulcanize/vulcanizedb/libraries/shared/transformer"

	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

func GetVatTuneConfig() shared_t.TransformerConfig {
	return shared_t.TransformerConfig{
		TransformerName:     constants.VatTuneLabel,
		ContractAddresses:   []string{constants.OldVatContractAddress()},
		ContractAbi:         constants.OldVatABI(),
		Topic:               constants.GetVatTuneSignature(),
		StartingBlockNumber: constants.OldVatDeploymentBlock(),
		EndingBlockNumber:   -1,
	}
}
