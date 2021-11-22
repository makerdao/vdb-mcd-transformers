package v2_0_0

import (
"github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip/initializers"
"github.com/makerdao/vdb-transformer-utilities/pkg/shared/constants"
)

var contractAddress = constants.GetContractAddress("MCD_CLIP_WBTC_B_2_0_0")
var StorageTransformerInitializer = initializers.GenerateStorageTransformerInitializer(contractAddress)