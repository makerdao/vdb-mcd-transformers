package initializer

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/mcd_transformers/transformers/storage/flap"

	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"

	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	s2 "github.com/vulcanize/vulcanizedb/libraries/shared/factories/storage"
)

var StorageTransformerInitializer transformer.StorageTransformerInitializer = s2.Transformer{
	Address: common.HexToAddress(constants.FlapperContractAddress()),
	Mappings: &flap.StorageKeysLookup{
		StorageRepository: &storage.MakerStorageRepository{},
		ContractAddress:   constants.FlapperContractAddress(),
	},
	Repository: &flap.FlapStorageRepository{ContractAddress: constants.FlapperContractAddress()},
}.NewTransformer
