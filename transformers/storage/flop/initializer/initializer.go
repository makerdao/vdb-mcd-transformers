package initializer

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/mcd_transformers/transformers/storage/flop"
	s2 "github.com/vulcanize/vulcanizedb/libraries/shared/factories/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
)

var StorageTransformerInitializer transformer.StorageTransformerInitializer = s2.Transformer{
	Address: common.HexToAddress(constants.FlopperContractAddress()),
	Mappings: &flop.StorageKeysLookup{
		StorageRepository: &storage.MakerStorageRepository{},
		ContractAddress:   constants.FlopperContractAddress()},
	Repository: &flop.FlopStorageRepository{ContractAddress: constants.FlopperContractAddress()},
}.NewTransformer
