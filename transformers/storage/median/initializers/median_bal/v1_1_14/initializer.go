// This is a plugin generated to export the configured transformer initializers

package v1_1_14

import (
	constants "github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	initializers "github.com/makerdao/vdb-mcd-transformers/transformers/storage/median/initializers"
)

var contractAddress = constants.GetContractAddress("MEDIAN_BAL_1_1_14")
var StorageTransformerInitializer = initializers.GenerateStorageTransformerInitializer(contractAddress)
