package cat

import (
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/libraries/shared/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage_diffs/maker"
)

const (
	NFlip = "nflip"
	Live  = "live"
	Vat   = "vat"
	Pit   = "pit"
	Vow   = "vow"

	IlkFlip = "flip"
	IlkChop = "chop"
	IlkLump = "lump"

	FlipIlk = "ilk"
	FlipUrn = "urn"
	FlipInk = "ink"
	FlipTab = "tab"
)

var (
	// wards takes up index 0
	IlksMappingIndex  = storage.IndexOne // bytes32 => flip address; chop (ray), lump (wad) uint256
	FlipsMappingIndex = storage.IndexTwo // uint256 => ilk, urn bytes32; ink, tab uint256 (both wad)

	NFlipKey      = common.HexToHash(storage.IndexThree)
	NFlipMetadata = utils.GetStorageValueMetadata(NFlip, nil, utils.Uint256)

	LiveKey      = common.HexToHash(storage.IndexFour)
	LiveMetadata = utils.GetStorageValueMetadata(Live, nil, utils.Uint256)

	VatKey      = common.HexToHash(storage.IndexFive)
	VatMetadata = utils.GetStorageValueMetadata(Vat, nil, utils.Address)

	PitKey      = common.HexToHash(storage.IndexSix)
	PitMetadata = utils.GetStorageValueMetadata(Pit, nil, utils.Address)

	VowKey      = common.HexToHash(storage.IndexSeven)
	VowMetadata = utils.GetStorageValueMetadata(Vow, nil, utils.Address)
)

type CatMappings struct {
	StorageRepository maker.IMakerStorageRepository
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

	flipsErr := mappings.loadFlipsKeys()
	if flipsErr != nil {
		return flipsErr
	}

	return nil
}

func loadStaticMappings() map[common.Hash]utils.StorageValueMetadata {
	mappings := make(map[common.Hash]utils.StorageValueMetadata)
	mappings[NFlipKey] = NFlipMetadata
	mappings[LiveKey] = LiveMetadata
	mappings[VatKey] = VatMetadata
	mappings[PitKey] = PitMetadata
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

// Flip ID increments each time it happens, so we just need the biggest flip ID from the DB
// and we can interpolate the sequence [0..max]. This makes sure we track all earlier flips,
// even if we've missed events
func (mappings CatMappings) loadFlipsKeys() error {
	maxFlip, err := mappings.StorageRepository.GetMaxFlip()
	if err != nil {
		logrus.Error("loadFlipsKeys: error getting max flip: ", err)
		return err
	} else if maxFlip == nil { // No flips occurred yet
		return nil
	}

	last := maxFlip.Int64()
	for flip := 0; int64(flip) <= last; flip++ {
		flipStr := strconv.Itoa(flip)
		mappings.mappings[getFlipIlkKey(flipStr)] = getFlipIlkMetadata(flipStr)
		mappings.mappings[getFlipUrnKey(flipStr)] = getFlipUrnMetadata(flipStr)
		mappings.mappings[getFlipInkKey(flipStr)] = getFlipInkMetadata(flipStr)
		mappings.mappings[getFlipTabKey(flipStr)] = getFlipTabMetadata(flipStr)
	}
	return nil
}

func getFlipIlkKey(flip string) common.Hash {
	return storage.GetMapping(FlipsMappingIndex, flip)
}

func getFlipIlkMetadata(flip string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Flip: flip}
	return utils.GetStorageValueMetadata(FlipIlk, keys, utils.Bytes32)
}

func getFlipUrnKey(flip string) common.Hash {
	return storage.GetIncrementedKey(getFlipIlkKey(flip), 1)
}

func getFlipUrnMetadata(flip string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Flip: flip}
	return utils.GetStorageValueMetadata(FlipUrn, keys, utils.Bytes32)
}

func getFlipInkKey(flip string) common.Hash {
	return storage.GetIncrementedKey(getFlipIlkKey(flip), 2)
}

func getFlipInkMetadata(flip string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Flip: flip}
	return utils.GetStorageValueMetadata(FlipInk, keys, utils.Uint256)
}

func getFlipTabKey(flip string) common.Hash {
	return storage.GetIncrementedKey(getFlipIlkKey(flip), 3)
}

func getFlipTabMetadata(flip string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Flip: flip}
	return utils.GetStorageValueMetadata(FlipTab, keys, utils.Uint256)
}
