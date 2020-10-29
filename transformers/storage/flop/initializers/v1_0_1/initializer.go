package v1_0_1

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/flop/initializers"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared/constants"
)

var StorageTransformerInitializer = initializers.GenerateStorageTransformerInitializer(constants.GetContractAddress("MCD_FLOP_1_0_1"))
