package initializer

import (
	"github.com/ethereum/go-ethereum/common"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/cat"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/cat/v1_1_0"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
)

var catAddress = constants.GetContractAddress("MCD_CAT_1_1_0")
var StorageTransformerInitializer storage.TransformerInitializer = storage.Transformer{
	Address:           common.HexToAddress(catAddress),
	StorageKeysLookup: storage.NewKeysLookup(v1_1_0.NewKeysLoader(&mcdStorage.MakerStorageRepository{}, catAddress)),
	Repository:        &cat.StorageRepository{ContractAddress: catAddress},
}.NewTransformer
