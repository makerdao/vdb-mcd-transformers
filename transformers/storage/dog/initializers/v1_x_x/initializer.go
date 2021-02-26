package v1_x_x

import (
	"github.com/ethereum/go-ethereum/common"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/cat"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/dog"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
)

var dogAddress = constants.GetContractAddress("MCD_Dog_1_x_x")
var StorageTransformerInitializer storage.TransformerInitializer = storage.Transformer{
	Address:           common.HexToAddress(dogAddress),
	StorageKeysLookup: storage.NewKeysLookup(dog.NewKeysLoader(&mcdStorage.MakerStorageRepository{}, dogAddress)),
	Repository:        &cat.StorageRepository{ContractAddress: dogAddress},
}.NewTransformer

