package dog

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities/wards"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	vdbStorage "github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

var (
	//wards at index0
	VatStorageKey = common.HexToHash(vdbStorage.IndexOne)
	VatMetadata   = types.GetValueMetadata(Vat, nil, types.Address)
	//ilks at index2
	VowStorageKey = common.HexToHash(vdbStorage.IndexThree)
	VowMetadata   = types.GetValueMetadata(Vow, nil, types.Address)

	LiveStorageKey = common.HexToHash(vdbStorage.IndexFour)
	LiveMetadata   = types.GetValueMetadata(Live, nil, types.Uint256)

	HoleStorageKey = common.HexToHash(vdbStorage.IndexFive)
	HoleMetadata   = types.GetValueMetadata(Hole, nil, types.Uint256)

	DirtStorageKey = common.HexToHash(vdbStorage.IndexSix)
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
		return nil, fmt.Errorf("error adding wards keys to cat keys loader: %w", wardsErr)
	}

	return mappings, nil
}

func loadStaticMappings() map[common.Hash]types.ValueMetadata {
	mappings := make(map[common.Hash]types.ValueMetadata)
	mappings[VatStorageKey] = VatMetadata
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
