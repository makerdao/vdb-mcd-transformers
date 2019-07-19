package initializer

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/vulcanize/vulcanizedb/libraries/shared/factories/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"

	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	storage2 "github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/mcd_transformers/transformers/storage/jug"
)

var StorageTransformerInitializer transformer.StorageTransformerInitializer = storage.Transformer{
	Address:    common.HexToAddress(constants.GetContractAddress("MCD_JUG")),
	Mappings:   &jug.JugMappings{StorageRepository: &storage2.MakerStorageRepository{}},
	Repository: &jug.JugStorageRepository{},
}.NewTransformer
