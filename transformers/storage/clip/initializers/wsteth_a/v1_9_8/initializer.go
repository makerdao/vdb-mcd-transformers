package v1_9_8

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip/initializers"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared/constants"
)

var contractAddress = constants.GetContractAddress("MCD_CLIP_WSTETH_A_1_9_8")
var StorageTransformerInitializer = initializers.GenerateStorageTransformerInitializer(contractAddress)
