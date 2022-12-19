package v12_16_22

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip/initializers"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared/constants"
)

var contractAddress = constants.GetContractAddress("MCD_CLIP_GNO_A_12_16_22")
var StorageTransformerInitializer = initializers.GenerateStorageTransformerInitializer(contractAddress)
