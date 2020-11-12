// This is a plugin generated to export the configured transformer initializers

package v1_1_14

import (
	initializers "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers"
	constants "github.com/makerdao/vdb-transformer-utilities/pkg/shared/constants"
)

var contractAddress = constants.GetContractAddress("MCD_FLIP_YFI_A_1_1_14")
var StorageTransformerInitializer = initializers.GenerateStorageTransformerInitializer(contractAddress)
