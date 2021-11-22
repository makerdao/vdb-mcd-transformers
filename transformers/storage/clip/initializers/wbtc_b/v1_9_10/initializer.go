package v1_9_10

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip/initializers"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared/constants"
)

var contractAddress = constants.GetContractAddress("MCD_CLIP_WBTC_B_1_9_10")
var StorageTransformerInitializer = initializers.GenerateStorageTransformerInitializer(contractAddress)
