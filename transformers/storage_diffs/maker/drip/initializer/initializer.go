package initializer

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage_diffs/maker"
	"github.com/vulcanize/mcd_transformers/transformers/storage_diffs/maker/drip"
	"github.com/vulcanize/vulcanizedb/libraries/shared/factories/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
)

var StorageTransformerInitializer transformer.StorageTransformerInitializer = storage.Transformer{
	Address:    common.HexToAddress(constants.DripContractAddress()),
	Mappings:   &drip.DripMappings{StorageRepository: &maker.MakerStorageRepository{}},
	Repository: &drip.DripStorageRepository{},
}.NewTransformer
