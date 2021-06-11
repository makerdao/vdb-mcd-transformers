package v1_6_0

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip/initializers"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared/constants"
)

var contractAddress = constants.GetContractAddress("MCD_CLIP_BAL_A_1_6_0")
var StorageTransformerInitializer = initializers.GenerateStorageTransformerInitializer(contractAddress)
