package usdc_a_flip

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers"
)

var StorageTransformerInitializer = initializers.GenerateStorageTransformerInitializer(constants.GetContractAddress("MCD_FLIP_USDC_A"))
