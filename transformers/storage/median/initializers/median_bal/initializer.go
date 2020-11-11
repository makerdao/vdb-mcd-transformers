// This is a plugin generated to export the configured transformer initializers

package median_bal

import (
	initializers "github.com/makerdao/vdb-mcd-transformers/transformers/storage/median/initializers"
	constants "github.com/makerdao/vdb-transformer-utilities/pkg/shared/constants"
)

var contractAddress = constants.GetContractAddress("MEDIAN_BAL")
var StorageTransformerInitializer = initializers.GenerateStorageTransformerInitializer(contractAddress)
