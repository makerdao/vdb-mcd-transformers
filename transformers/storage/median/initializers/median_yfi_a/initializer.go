// This is a plugin generated to export the configured transformer initializers

package median_yfi_a

import (
	initializers "github.com/makerdao/vdb-mcd-transformers/transformers/storage/median/initializers"
	constants "github.com/makerdao/vdb-transformer-utilities/pkg/shared/constants"
)

var contractAddress = constants.GetContractAddress("MEDIAN_YFI_A")
var StorageTransformerInitializer = initializers.GenerateStorageTransformerInitializer(contractAddress)
