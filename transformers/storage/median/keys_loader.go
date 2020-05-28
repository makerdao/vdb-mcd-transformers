package median

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
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
	Val  = "val"
	Age  = "age"
	Bar  = "bar"
	Bud  = "bud"
	Orcl = "orcl"
	Slot = "slot"
)

var (
	ValAndAgeStorageKey = common.HexToHash(vdbStorage.IndexOne)
	valAndAgeTypes      = map[int]types.ValueType{0: types.Uint128, 1: types.Uint32}
	valAndAgeNames      = map[int]string{0: Val, 1: Age}
	ValAndAgeMetadata   = types.GetValueMetadataForPackedSlot(mcdStorage.Packed, nil, types.PackedSlot, valAndAgeNames, valAndAgeTypes)

	BarKey      = common.HexToHash(vdbStorage.IndexTwo)
	BarMetadata = types.GetValueMetadata(Bar, nil, types.Uint256)

	OrclMappingIndex = vdbStorage.IndexThree
	BudMappingIndex  = vdbStorage.IndexFour
	SlotMappingIndex = vdbStorage.IndexFive
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
	return loader.addDynamicMappings(mappings)
}

func (loader keysLoader) addDynamicMappings(mappings map[common.Hash]types.ValueMetadata) (map[common.Hash]types.ValueMetadata, error) {
	mappings, wardsErr := loader.loadWardsKeys(mappings)
	if wardsErr != nil {
		return nil, wardsErr
	}
	mappings, orclErr := loader.loadOrclKeys(mappings)
	if orclErr != nil {
		return nil, orclErr
	}
	mappings, slotErr := loader.loadSlotKeys(mappings)
	if slotErr != nil {
		return nil, slotErr
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

func (loader *keysLoader) loadOrclKeys(mappings map[common.Hash]types.ValueMetadata) (map[common.Hash]types.ValueMetadata, error) {
	orclAddresses, orclErr := loader.storageRepository.GetMedianOrclAddresses(loader.contractAddress)
	if orclErr != nil {
		return nil, orclErr
	}
	for _, address := range orclAddresses {
		paddedAddress, padErr := utilities.PadAddress(address)
		if padErr != nil {
			return nil, padErr
		}
		mappings[getOrclKey(paddedAddress)] = getOrclMetadata(address)
	}
	return mappings, nil
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

func (loader *keysLoader) loadSlotKeys(mappings map[common.Hash]types.ValueMetadata) (map[common.Hash]types.ValueMetadata, error) {
	slotIds, slotErr := loader.storageRepository.GetMedianSlotIds()
	if slotErr != nil {
		return nil, slotErr
	}
	for _, slotId := range slotIds {
		hexSlotId, convertErr := shared.ConvertIntStringToHex(slotId)
		if convertErr != nil {
			return nil, convertErr
		}
		mappings[getSlotKey(hexSlotId)] = getSlotMetadata(slotId)
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

func getOrclKey(address string) common.Hash {
	return vdbStorage.GetKeyForMapping(OrclMappingIndex, address)
}

func getOrclMetadata(address string) types.ValueMetadata {
	keys := map[types.Key]string{constants.Address: address}
	return types.GetValueMetadata(Orcl, keys, types.Uint256)
}

func getSlotKey(slotId string) common.Hash {
	return vdbStorage.GetKeyForMapping(SlotMappingIndex, slotId)
}

func getSlotMetadata(slotId string) types.ValueMetadata {
	keys := map[types.Key]string{constants.SlotId: slotId}
	return types.GetValueMetadata(Slot, keys, types.Address)
}

func (loader keysLoader) SetDB(db *postgres.DB) {
	loader.storageRepository.SetDB(db)
}
