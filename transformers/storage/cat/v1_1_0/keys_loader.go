package v1_1_0

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/cat"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	vdbStorage "github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

var (
	BoxKey      = common.HexToHash(vdbStorage.IndexFive)
	BoxMetadata = types.GetValueMetadata(cat.Box, nil, types.Uint256)

	LitterKey      = common.HexToHash(vdbStorage.IndexSix)
	LitterMetadata = types.GetValueMetadata(cat.Litter, nil, types.Uint256)
)

type keysLoader struct {
	storageRepository mcdStorage.IMakerStorageRepository
	contractAddress   string
}

func NewKeysLoader(storageRepository mcdStorage.IMakerStorageRepository, contractAddress string) storage.KeysLoader {
	return &keysLoader{storageRepository: storageRepository, contractAddress: contractAddress}
}

func (loader *keysLoader) SetDB(db *postgres.DB) {
	loader.storageRepository.SetDB(db)
}

func (loader *keysLoader) LoadMappings() (map[common.Hash]types.ValueMetadata, error) {
	mappings := loadStaticMappings()
	mappings, sharedErr := cat.LoadSharedMappings(mappings, loader.contractAddress, loader.storageRepository)
	if sharedErr != nil {
		return nil, fmt.Errorf("error adding shared cat keys to v1_0_0 keys loader: %w", sharedErr)
	}
	mappings, ilkErr := loader.addIlkKeys(mappings)
	if ilkErr != nil {
		return nil, fmt.Errorf("error adding ilk keys to cat keys loader: %w", ilkErr)
	}
	return mappings, nil
}

func loadStaticMappings() map[common.Hash]types.ValueMetadata {
	mappings := make(map[common.Hash]types.ValueMetadata)
	mappings[BoxKey] = BoxMetadata
	mappings[LitterKey] = LitterMetadata
	return mappings
}

func (loader *keysLoader) addIlkKeys(mappings map[common.Hash]types.ValueMetadata) (map[common.Hash]types.ValueMetadata, error) {
	ilks, err := loader.storageRepository.GetIlks()
	if err != nil {
		return nil, fmt.Errorf("error getting ilks: %w", err)
	}
	for _, ilk := range ilks {
		mappings[cat.GetIlkFlipKey(ilk)] = cat.GetIlkFlipMetadata(ilk)
		mappings[cat.GetIlkChopKey(ilk)] = cat.GetIlkChopMetadata(ilk)
		mappings[getIlkDunkKey(ilk)] = getIlkDunkMetadata(ilk)
	}
	return mappings, nil
}

func getIlkDunkKey(ilk string) common.Hash {
	return vdbStorage.GetIncrementedKey(cat.GetIlkFlipKey(ilk), 2)
}

func getIlkDunkMetadata(ilk string) types.ValueMetadata {
	keys := map[types.Key]string{constants.Ilk: ilk}
	return types.GetValueMetadata(cat.IlkDunk, keys, types.Uint256)
}
