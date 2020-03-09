package medianizer

import (
	"github.com/ethereum/go-ethereum/common"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	vdbStorage "github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

const (
	Val = "val"
	Age = "age"
	Bar = "bar"
)

var (
	ValKey      = common.HexToHash(vdbStorage.IndexOne)
	ValMetadata = types.GetValueMetadata(Val, nil, types.Uint128)

	AgeKey      = common.HexToHash(vdbStorage.IndexTwo)
	AgeMetadata = types.GetValueMetadata(Age, nil, types.Bytes32)

	BarKey      = common.HexToHash(vdbStorage.IndexFour)
	BarMetadata = types.GetValueMetadata(Bar, nil, types.Uint256)
)

type keysLoader struct {
	storageRepository mcdStorage.IMakerStorageRepository
	contractAddress   string
}

func (loader keysLoader) LoadMappings() (map[common.Hash]types.ValueMetadata, error) {
	mappings := make(map[common.Hash]types.ValueMetadata)
	mappings[ValKey] = ValMetadata
	mappings[AgeKey] = AgeMetadata
	mappings[BarKey] = BarMetadata
	return mappings, nil
}

func (loader keysLoader) SetDB(db *postgres.DB) {
	loader.storageRepository.SetDB(db)
}

func NewKeysLoader(storageRepository mcdStorage.IMakerStorageRepository, contractAddress string) storage.KeysLoader {
	return &keysLoader{storageRepository: storageRepository, contractAddress: contractAddress}
}
