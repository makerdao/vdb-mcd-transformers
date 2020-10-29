package v1_0_9

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/flap/initializers"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared/constants"
)

var StorageTransformerInitializer = initializers.GenerateStorageTransformerInitializer(constants.GetContractAddress("MCD_FLAP_1_0_9"))
