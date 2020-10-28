package v1_0_9

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/flop/initializers"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared/constants"
)

var StorageTransformerInitializer = initializers.GenerateStorageTransformerInitializer(constants.GetContractAddress("MCD_FLOP_1.0.9"))
