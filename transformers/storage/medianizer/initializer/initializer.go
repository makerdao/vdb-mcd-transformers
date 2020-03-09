package initializer

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/medianizer"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
)

var medianizerAddress = constants.GetContractAddress("MCD_MEDIANIZER")
var StorageTransformerInitializer storage.TransformerInitializer = storage.Transformer{
	Address:           common.HexToAddress(medianizerAddress),
	StorageKeysLookup: storage.NewKeysLookup(medianizer.NewKeysLoader(&mcdStorage.MakerStorageRepository{}, medianizerAddress)),
	Repository:        &medianizer.MedianizerStorageRepository{ContractAddress: medianizerAddress},
}.NewTransformer
