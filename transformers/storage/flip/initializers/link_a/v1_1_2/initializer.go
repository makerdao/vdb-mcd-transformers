package v1_1_2

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared/constants"
)

var StorageTransformerInitializer = initializers.GenerateStorageTransformerInitializer(constants.GetContractAddress("MCD_FLIP_LINK_A_1.1.2"))
