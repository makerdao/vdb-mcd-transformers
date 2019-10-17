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

type StorageKeysLookup struct {
	StorageRepository s2.IMakerStorageRepository
	mappings          map[common.Hash]utils.StorageValueMetadata
}

func (lookup StorageKeysLookup) Lookup(key common.Hash) (utils.StorageValueMetadata, error) {
	metadata, ok := lookup.mappings[key]
	if !ok {
		err := lookup.loadMappings()
		if err != nil {
			return metadata, err
		}
		metadata, ok = lookup.mappings[key]
		if !ok {
			return metadata, utils.ErrStorageKeyNotFound{Key: key.Hex()}
		}
	}
	return metadata, nil
}

func (lookup *StorageKeysLookup) SetDB(db *postgres.DB) {
	lookup.StorageRepository.SetDB(db)
}

func (lookup *StorageKeysLookup) loadMappings() error {
	lookup.mappings = loadStaticMappings()
	ilkErr := lookup.loadIlkKeys()
	if ilkErr != nil {
		return ilkErr
	}
	lookup.mappings = storage.AddHashedKeys(lookup.mappings)
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
func (lookup *StorageKeysLookup) loadIlkKeys() error {
	ilks, err := lookup.StorageRepository.GetIlks()
	if err != nil {
		return err
	}
	for _, ilk := range ilks {
		lookup.mappings[getIlkFlipKey(ilk)] = getIlkFlipMetadata(ilk)
		lookup.mappings[getIlkChopKey(ilk)] = getIlkChopMetadata(ilk)
		lookup.mappings[getIlkLumpKey(ilk)] = getIlkLumpMetadata(ilk)
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
