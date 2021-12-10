package v1_9_12

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip/initializers"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared/constants"
)

var contractAddress = constants.GetContractAddress("MCD_CLIP_GUNIV3DAIUSDC_2_1_9_12")
var StorageTransformerInitializer = initializers.GenerateStorageTransformerInitializer(contractAddress)
