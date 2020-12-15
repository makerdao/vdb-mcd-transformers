// This is a plugin generated to export the configured transformer initializers

package v1_2_1

import (
	constants "github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	initializers "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers"
)

var contractAddress = constants.GetContractAddress("MCD_FLIP_RENBTC_A_1_2_1")
var StorageTransformerInitializer = initializers.GenerateStorageTransformerInitializer(contractAddress)
