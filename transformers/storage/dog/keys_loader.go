package dog

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities/wards"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	vdbStorage "github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

var (
	//wards at index0
	IlksMappingIndex = vdbStorage.IndexOne // bytes32 => clip address; chop (wad), hole (rad), dirt (rad) uint256

	VowStorageKey = common.HexToHash(vdbStorage.IndexTwo)
	VowMetadata   = types.GetValueMetadata(Vow, nil, types.Address)

	LiveStorageKey = common.HexToHash(vdbStorage.IndexThree)
	LiveMetadata   = types.GetValueMetadata(Live, nil, types.Uint256)

	HoleStorageKey = common.HexToHash(vdbStorage.IndexFour)
	HoleMetadata   = types.GetValueMetadata(Hole, nil, types.Uint256)

	DirtStorageKey = common.HexToHash(vdbStorage.IndexFive)
	DirtMetadata   = types.GetValueMetadata(Dirt, nil, types.Uint256)
)

type keysLoader struct {
	storageRepository mcdStorage.IMakerStorageRepository
	contractAddress   string
}

func NewKeysLoader(storageRepository mcdStorage.IMakerStorageRepository, contractAddress string) storage.KeysLoader {
	return &keysLoader{
		storageRepository: storageRepository,
		contractAddress:   contractAddress,
	}
}

func (loader *keysLoader) SetDB(db *postgres.DB) {
	loader.storageRepository.SetDB(db)
}

func (loader *keysLoader) LoadMappings() (map[common.Hash]types.ValueMetadata, error) {
	mappings := loadStaticMappings()
	mappings, wardsErr := addWardsKeys(mappings, loader.contractAddress, loader.storageRepository)
	if wardsErr != nil {
		return nil, fmt.Errorf("error adding wards keys to dog keys loader: %w", wardsErr)
	}
	mappings, ilkErr := loader.addIlkKeys(mappings)
	if ilkErr != nil {
		return nil, fmt.Errorf("error adding ilk keys to dog keys loader: %w", ilkErr)
	}
	return mappings, nil
}

func loadStaticMappings() map[common.Hash]types.ValueMetadata {
	mappings := make(map[common.Hash]types.ValueMetadata)
	mappings[VowStorageKey] = VowMetadata
	mappings[LiveStorageKey] = LiveMetadata
	mappings[HoleStorageKey] = HoleMetadata
	mappings[DirtStorageKey] = DirtMetadata
	return mappings
}

func addWardsKeys(mappings map[common.Hash]types.ValueMetadata, address string, repository mcdStorage.IMakerStorageRepository) (map[common.Hash]types.ValueMetadata, error) {
	addresses, err := repository.GetWardsAddresses(address)
	if err != nil {
		return nil, fmt.Errorf("error getting wards addresses: %w", err)
	}
	return wards.AddWardsKeys(mappings, addresses)
}

func (loader *keysLoader) addIlkKeys(mappings map[common.Hash]types.ValueMetadata) (map[common.Hash]types.ValueMetadata, error) {
	ilks, err := loader.storageRepository.GetIlks()
	if err != nil {
		return nil, fmt.Errorf("error getting ilks: %w", err)
	}
	for _, ilk := range ilks {
		mappings[GetIlkClipKey(ilk)] = GetIlkClipMetadata(ilk)
		mappings[GetIlkChopKey(ilk)] = GetIlkChopMetadata(ilk)
		mappings[GetIlkHoleKey(ilk)] = GetIlkHoleMetadata(ilk)
		mappings[GetIlkDirtKey(ilk)] = GetIlkDirtMetadata(ilk)
	}
	return mappings, nil
}

func GetIlkClipKey(ilk string) common.Hash {
	return vdbStorage.GetKeyForMapping(IlksMappingIndex, ilk)
}

func GetIlkClipMetadata(ilk string) types.ValueMetadata {
	keys := map[types.Key]string{constants.Ilk: ilk}
	return types.GetValueMetadata(IlkClip, keys, types.Address)
}

func GetIlkChopKey(ilk string) common.Hash {
	return vdbStorage.GetIncrementedKey(GetIlkClipKey(ilk), 1)
}

func GetIlkChopMetadata(ilk string) types.ValueMetadata {
	keys := map[types.Key]string{constants.Ilk: ilk}
	return types.GetValueMetadata(IlkChop, keys, types.Uint256)
}

func GetIlkHoleKey(ilk string) common.Hash {
	return vdbStorage.GetIncrementedKey(GetIlkClipKey(ilk), 2)
}

func GetIlkHoleMetadata(ilk string) types.ValueMetadata {
	keys := map[types.Key]string{constants.Ilk: ilk}
	return types.GetValueMetadata(IlkHole, keys, types.Uint256)
}

func GetIlkDirtKey(ilk string) common.Hash {
	return vdbStorage.GetIncrementedKey(GetIlkClipKey(ilk), 3)
}

func GetIlkDirtMetadata(ilk string) types.ValueMetadata {
	keys := map[types.Key]string{constants.Ilk: ilk}
	return types.GetValueMetadata(IlkDirt, keys, types.Uint256)
}
