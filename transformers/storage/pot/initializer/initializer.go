package initializer

import (
	"github.com/ethereum/go-ethereum/common"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/pot"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
)

var potAddress = constants.GetContractAddress("MCD_POT")
var StorageTransformerInitializer storage.TransformerInitializer = storage.Transformer{
	Address:           common.HexToAddress(constants.GetContractAddress("MCD_POT")),
	StorageKeysLookup: storage.NewKeysLookup(pot.NewKeysLoader(&mcdStorage.MakerStorageRepository{}, potAddress)),
	Repository:        &pot.StorageRepository{},
}.NewTransformer
