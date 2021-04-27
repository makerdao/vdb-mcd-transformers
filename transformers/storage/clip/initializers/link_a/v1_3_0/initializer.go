package v1_3_0

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip/initializers"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared/constants"
)

var contractAddress = constants.GetContractAddress("MCD_CLIP_LINK_A_1_3_0")
var StorageTransformerInitializer = initializers.GenerateStorageTransformerInitializer(contractAddress)
