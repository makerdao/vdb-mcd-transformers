package median

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
	ValAndAgeStorageKey = common.HexToHash(vdbStorage.IndexOne)
	valAndAgeTypes      = map[int]types.ValueType{0: types.Uint128, 1: types.Uint32}
	valAndAgeNames      = map[int]string{0: Val, 1: Age}
	ValAndAgeMetadata   = types.GetValueMetadataForPackedSlot(mcdStorage.Packed, nil, types.PackedSlot, valAndAgeNames, valAndAgeTypes)
	BarKey              = common.HexToHash(vdbStorage.IndexTwo)
	BarMetadata         = types.GetValueMetadata(Bar, nil, types.Uint256)
)

type keysLoader struct {
	storageRepository mcdStorage.IMakerStorageRepository
	contractAddress   string
}

func NewKeysLoader(storageRepository mcdStorage.IMakerStorageRepository, contractAddress string) storage.KeysLoader {
	return &keysLoader{storageRepository: storageRepository, contractAddress: contractAddress}
}

func (loader keysLoader) LoadMappings() (map[common.Hash]types.ValueMetadata, error) {
	mappings := make(map[common.Hash]types.ValueMetadata)
	mappings[ValAndAgeStorageKey] = ValAndAgeMetadata
	mappings[BarKey] = BarMetadata
	return mappings, nil
}

func (loader keysLoader) SetDB(db *postgres.DB) {
	loader.storageRepository.SetDB(db)
}
