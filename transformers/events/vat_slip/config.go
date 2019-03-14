package vat_slip

import (
	shared_t "github.com/vulcanize/vulcanizedb/libraries/shared/transformer"

	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

func GetVatSlipConfig() shared_t.TransformerConfig {
	return shared_t.TransformerConfig{
		TransformerName:     constants.VatSlipLabel,
		ContractAddresses:   []string{constants.OldVatContractAddress()},
		ContractAbi:         constants.OldVatABI(),
		Topic:               constants.GetVatSlipSignature(),
		StartingBlockNumber: constants.OldVatDeploymentBlock(),
		EndingBlockNumber:   -1,
	}
}
