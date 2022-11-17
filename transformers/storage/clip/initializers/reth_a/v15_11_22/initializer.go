package v15_11_22

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip/initializers"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared/constants"
)

var contractAddress = constants.GetContractAddress("MCD_CLIP_RETH_A_15_11_22")
var StorageTransformerInitializer = initializers.GenerateStorageTransformerInitializer(contractAddress)
