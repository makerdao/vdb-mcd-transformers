package v1_0_9

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared/constants"
)

var StorageTransformerInitializer = initializers.GenerateStorageTransformerInitializer(constants.GetContractAddress("MCD_FLIP_USDC_A_1.0.9"))
