package initializer

import (
	"github.com/ethereum/go-ethereum/common"

	s2 "github.com/vulcanize/vulcanizedb/libraries/shared/factories/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"

	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/mcd_transformers/transformers/storage/flop"
)

var StorageTransformerInitializer transformer.StorageTransformerInitializer = s2.Transformer{
	Address: common.HexToAddress(constants.GetContractAddress("MCD_FLOP")),
	Mappings: &flop.StorageKeysLookup{
		StorageRepository: &storage.MakerStorageRepository{},
		ContractAddress:   constants.GetContractAddress("MCD_FLOP")},
	Repository: &flop.FlopStorageRepository{ContractAddress: constants.GetContractAddress("MCD_FLOP")},
}.NewTransformer
