package v1_11_0

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip/initializers"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared/constants"
)

var contractAddress = constants.GetContractAddress("MCD_CLIP_GUNIV3DAIUSDC_1_1_9_6")
var StorageTransformerInitializer = initializers.GenerateStorageTransformerInitializer(contractAddress)
