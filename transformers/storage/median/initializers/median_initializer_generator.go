package initializers

import (
	"github.com/ethereum/go-ethereum/common"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/median"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
)

func GenerateStorageTransformerInitializer(contractAddress string) storage.TransformerInitializer {
	return storage.Transformer{
		Address:           common.HexToAddress(contractAddress),
		StorageKeysLookup: storage.NewKeysLookup(median.NewKeysLoader(&mcdStorage.MakerStorageRepository{}, contractAddress)),
		Repository:        &median.MedianizerStorageRepository{ContractAddress: contractAddress},
	}.NewTransformer
}
