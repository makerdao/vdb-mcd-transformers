package initializer

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/vulcanize/vulcanizedb/libraries/shared/factories/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"

	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	storage2 "github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/mcd_transformers/transformers/storage/drip"
)

var StorageTransformerInitializer transformer.StorageTransformerInitializer = storage.Transformer{
	Address:    common.HexToAddress(constants.DripContractAddress()),
	Mappings:   &drip.DripMappings{StorageRepository: &storage2.MakerStorageRepository{}},
	Repository: &drip.DripStorageRepository{},
}.NewTransformer
