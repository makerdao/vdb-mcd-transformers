package v1_9_4

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip/initializers"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared/constants"
)

var contractAddress = constants.GetContractAddress("MCD_CLIP_MATIC_A_1_9_4")
var StorageTransformerInitializer = initializers.GenerateStorageTransformerInitializer(contractAddress)
