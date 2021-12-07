package v1_9_11

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip/initializers"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared/constants"
)

var contractAddress = constants.GetContractAddress("MCD_CLIP_WBTC_C_1_9_11")
var StorageTransformerInitializer = initializers.GenerateStorageTransformerInitializer(contractAddress)
