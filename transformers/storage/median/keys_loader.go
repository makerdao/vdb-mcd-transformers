package median

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities"
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
	Bud = "bud"
)

var (
	ValAndAgeStorageKey = common.HexToHash(vdbStorage.IndexOne)
	valAndAgeTypes      = map[int]types.ValueType{0: types.Uint128, 1: types.Uint32}
	valAndAgeNames      = map[int]string{0: Val, 1: Age}
	ValAndAgeMetadata   = types.GetValueMetadataForPackedSlot(mcdStorage.Packed, nil, types.PackedSlot, valAndAgeNames, valAndAgeTypes)

	BarKey      = common.HexToHash(vdbStorage.IndexTwo)
	BarMetadata = types.GetValueMetadata(Bar, nil, types.Uint256)

	BudMappingIndex = vdbStorage.IndexFour
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
	mappings, wardsErr := loader.loadWardsKeys(mappings)
	if wardsErr != nil {
		return nil, wardsErr
	}
	return loader.loadBudKeys(mappings)
}

func loadStaticKeys() map[common.Hash]types.ValueMetadata {
	mappings := make(map[common.Hash]types.ValueMetadata)
	mappings[ValAndAgeStorageKey] = ValAndAgeMetadata
	mappings[BarKey] = BarMetadata
	return mappings
}

func (loader *keysLoader) loadWardsKeys(mappings map[common.Hash]types.ValueMetadata) (map[common.Hash]types.ValueMetadata, error) {
	addresses, err := loader.storageRepository.GetWardsAddresses(loader.contractAddress)
	if err != nil {
		return nil, err
	}
	return wards.AddWardsKeys(mappings, addresses)
}

func (loader *keysLoader) loadBudKeys(mappings map[common.Hash]types.ValueMetadata) (map[common.Hash]types.ValueMetadata, error) {
	budAddresses, budErr := loader.storageRepository.GetMedianBudAddresses(loader.contractAddress)
	if budErr != nil {
		return nil, budErr
	}
	for _, address := range budAddresses {
		paddedAddress, padErr := utilities.PadAddress(address)
		if padErr != nil {
			return nil, padErr
		}
		mappings[getBudKey(paddedAddress)] = getBudMetadata(address)
	}
	return mappings, nil
}

func getBudKey(address string) common.Hash {
	return vdbStorage.GetKeyForMapping(BudMappingIndex, address)
}

func getBudMetadata(address string) types.ValueMetadata {
	keys := map[types.Key]string{constants.A: address}
	return types.GetValueMetadata(Bud, keys, types.Uint256)
}

func (loader keysLoader) SetDB(db *postgres.DB) {
	loader.storageRepository.SetDB(db)
}
