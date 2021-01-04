// This is a plugin generated to export the configured transformer initializers

package v1_2_1

import (
	initializers "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers"
	constants "github.com/makerdao/vdb-transformer-utilities/pkg/shared/constants"
)

var contractAddress = constants.GetContractAddress("MCD_FLIP_UNI_A_1_2_1")
var StorageTransformerInitializer = initializers.GenerateStorageTransformerInitializer(contractAddress)
