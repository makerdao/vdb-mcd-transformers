package initializer

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/pot"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/libraries/shared/transformer"
)

var potAddress = constants.GetContractAddress("MCD_POT")
var StorageTransformerInitializer transformer.StorageTransformerInitializer = storage.Transformer{
	Address:           common.HexToAddress(constants.GetContractAddress("MCD_POT")),
	HashedAddress:     types.HexToKeccak256Hash(constants.GetContractAddress("MCD_POT")),
	StorageKeysLookup: storage.NewKeysLookup(pot.NewKeysLoader(&mcdStorage.MakerStorageRepository{}, potAddress)),
	Repository:        &pot.PotStorageRepository{},
}.NewTransformer
