package initializers

import (
	"github.com/ethereum/go-ethereum/common"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/flap"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
)

func GenerateStorageTransformerInitializer(contractAddress string) storage.TransformerInitializer {
	return storage.Transformer{
		Address:           common.HexToAddress(contractAddress),
		StorageKeysLookup: storage.NewKeysLookup(flap.NewKeysLoader(&mcdStorage.MakerStorageRepository{}, contractAddress)),
		Repository:        &flap.StorageRepository{ContractAddress: contractAddress},
	}.NewTransformer
}
