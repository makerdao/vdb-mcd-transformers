package median

import (
	"github.com/ethereum/go-ethereum/common"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities/wards"
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

	BarKey      = common.HexToHash(vdbStorage.IndexTwo)
	BarMetadata = types.GetValueMetadata(Bar, nil, types.Uint256)
)

type keysLoader struct {
	storageRepository mcdStorage.IMakerStorageRepository
	contractAddress   string
}

func NewKeysLoader(storageRepository mcdStorage.IMakerStorageRepository, contractAddress string) storage.KeysLoader {
	return &keysLoader{storageRepository: storageRepository, contractAddress: contractAddress}
}

func (loader keysLoader) LoadMappings() (map[common.Hash]types.ValueMetadata, error) {
	mappings := loadStaticKeys()
	mappings, wardsErr := loader.addWardsKeys(mappings)
	if wardsErr != nil {
		return nil, wardsErr
	}
	return mappings, nil
}

func loadStaticKeys() map[common.Hash]types.ValueMetadata {
	mappings := make(map[common.Hash]types.ValueMetadata)
	mappings[ValAndAgeStorageKey] = ValAndAgeMetadata
	mappings[BarKey] = BarMetadata
	return mappings
}

func (loader *keysLoader) addWardsKeys(mappings map[common.Hash]types.ValueMetadata) (map[common.Hash]types.ValueMetadata, error) {
	addresses, err := loader.storageRepository.GetWardsAddresses(loader.contractAddress)
	if err != nil {
		return nil, err
	}
	return wards.AddWardsKeys(mappings, addresses)
}

func (loader keysLoader) SetDB(db *postgres.DB) {
	loader.storageRepository.SetDB(db)
}
