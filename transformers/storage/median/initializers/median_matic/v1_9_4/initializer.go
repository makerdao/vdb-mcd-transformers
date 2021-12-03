package v1_9_4

import (
	initializers "github.com/makerdao/vdb-mcd-transformers/transformers/storage/median/initializers"
	constants "github.com/makerdao/vdb-transformer-utilities/pkg/shared/constants"
)

var contractAddress = constants.GetContractAddress("MEDIAN_MATIC_1_9_4")
var StorageTransformerInitializer = initializers.GenerateStorageTransformerInitializer(contractAddress)
