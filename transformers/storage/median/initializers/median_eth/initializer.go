package median_eth

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/median/initializers"
)

var StorageTransformerInitializer = initializers.GenerateStorageTransformerInitializer(constants.GetContractAddress("MEDIAN_ETH"))
