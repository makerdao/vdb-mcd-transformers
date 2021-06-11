package v1_8_0

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip/initializers"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared/constants"
)

var contractAddress = constants.GetContractAddress("MCD_CLIP_UNIV2UNIETH_A_1_8_0")
var StorageTransformerInitializer = initializers.GenerateStorageTransformerInitializer(contractAddress)
