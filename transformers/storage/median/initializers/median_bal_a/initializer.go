// This is a plugin generated to export the configured transformer initializers

package median_bal_a

import (
	constants "github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	initializers "github.com/makerdao/vdb-mcd-transformers/transformers/storage/median/initializers"
)

var contractAddress = constants.GetContractAddress("MEDIAN_BAL_A")
var StorageTransformerInitializer = initializers.GenerateStorageTransformerInitializer(contractAddress)
