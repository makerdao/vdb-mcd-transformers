package v1_0_6

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/median/initializers"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared/constants"
)

var StorageTransformerInitializer = initializers.GenerateStorageTransformerInitializer(constants.GetContractAddress("MEDIAN_WBTC_1_0_6"))
