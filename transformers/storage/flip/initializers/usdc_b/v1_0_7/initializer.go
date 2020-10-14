package v1_0_7

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared/constants"
)

var StorageTransformerInitializer = initializers.GenerateStorageTransformerInitializer(constants.GetContractAddress("MCD_FLIP_USDC_B_1_0_7"))
