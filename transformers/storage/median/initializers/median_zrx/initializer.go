package median_zrx

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/median/initializers"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared/constants"
)

var StorageTransformerInitializer = initializers.GenerateStorageTransformerInitializer(constants.GetContractAddress("MEDIAN_ZRX"))
