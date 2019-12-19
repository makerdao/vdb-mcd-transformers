package initializer

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/pot"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	vdbStorage "github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/transformer"
)

var StorageTransformerInitializer transformer.StorageTransformerInitializer = storage.Transformer{
	HashedAddress:     vdbStorage.HexToKeccak256Hash(constants.GetContractAddress("MCD_POT")),
	StorageKeysLookup: storage.NewKeysLookup(pot.NewKeysLoader(&mcdStorage.MakerStorageRepository{})),
	Repository:        &pot.PotStorageRepository{},
}.NewTransformer
