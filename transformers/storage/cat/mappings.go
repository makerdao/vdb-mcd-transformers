package cat

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	s2 "github.com/vulcanize/mcd_transformers/transformers/storage"
)

const (
	Live = "live"
	Vat  = "vat"
	Vow  = "vow"

	IlkFlip = "flip"
	IlkChop = "chop"
	IlkLump = "lump"
)

var (
	// wards takes up index 0
	IlksMappingIndex = storage.IndexOne // bytes32 => flip address; chop (ray), lump (wad) uint256

	LiveKey      = common.HexToHash(storage.IndexTwo)
	LiveMetadata = utils.GetStorageValueMetadata(Live, nil, utils.Uint256)

	VatKey      = common.HexToHash(storage.IndexThree)
	VatMetadata = utils.GetStorageValueMetadata(Vat, nil, utils.Address)

	VowKey      = common.HexToHash(storage.IndexFour)
	VowMetadata = utils.GetStorageValueMetadata(Vow, nil, utils.Address)
)

type CatMappings struct {
	StorageRepository s2.IMakerStorageRepository
	mappings          map[common.Hash]utils.StorageValueMetadata
}

func (mappings CatMappings) Lookup(key common.Hash) (utils.StorageValueMetadata, error) {
	metadata, ok := mappings.mappings[key]
	if !ok {
		err := mappings.loadMappings()
		if err != nil {
			return metadata, err
		}
		metadata, ok = mappings.mappings[key]
		if !ok {
			return metadata, utils.ErrStorageKeyNotFound{Key: key.Hex()}
		}
	}
	return metadata, nil
}

func (mappings *CatMappings) SetDB(db *postgres.DB) {
	mappings.StorageRepository.SetDB(db)
}

func (mappings *CatMappings) loadMappings() error {
	mappings.mappings = loadStaticMappings()
	ilkErr := mappings.loadIlkKeys()
	if ilkErr != nil {
		return ilkErr
	}
	mappings.mappings = storage.AddHashedKeys(mappings.mappings)
	return nil
}

func loadStaticMappings() map[common.Hash]utils.StorageValueMetadata {
	mappings := make(map[common.Hash]utils.StorageValueMetadata)
	mappings[LiveKey] = LiveMetadata
	mappings[VatKey] = VatMetadata
	mappings[VowKey] = VowMetadata
	return mappings
}

// Ilks
func (mappings *CatMappings) loadIlkKeys() error {
	ilks, err := mappings.StorageRepository.GetIlks()
	if err != nil {
		return err
	}
	for _, ilk := range ilks {
		mappings.mappings[getIlkFlipKey(ilk)] = getIlkFlipMetadata(ilk)
		mappings.mappings[getIlkChopKey(ilk)] = getIlkChopMetadata(ilk)
		mappings.mappings[getIlkLumpKey(ilk)] = getIlkLumpMetadata(ilk)
	}
	return nil
}

func getIlkFlipKey(ilk string) common.Hash {
	return storage.GetMapping(IlksMappingIndex, ilk)
}

func getIlkFlipMetadata(ilk string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Ilk: ilk}
	return utils.GetStorageValueMetadata(IlkFlip, keys, utils.Address)
}

func getIlkChopKey(ilk string) common.Hash {
	return storage.GetIncrementedKey(getIlkFlipKey(ilk), 1)
}

func getIlkChopMetadata(ilk string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Ilk: ilk}
	return utils.GetStorageValueMetadata(IlkChop, keys, utils.Uint256)
}

func getIlkLumpKey(ilk string) common.Hash {
	return storage.GetIncrementedKey(getIlkFlipKey(ilk), 2)
}

func getIlkLumpMetadata(ilk string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Ilk: ilk}
	return utils.GetStorageValueMetadata(IlkLump, keys, utils.Uint256)
}
