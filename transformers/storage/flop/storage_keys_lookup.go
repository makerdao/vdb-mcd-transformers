package flop

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	s2 "github.com/vulcanize/mcd_transformers/transformers/storage"
)

const (
	Vat    = "vat"
	Gem    = "gem"
	Beg    = "beg"
	Ttl    = "ttl"
	Tau    = "tau"
	Kicks  = "kicks"
	Live   = "live"
	Packed = "packed_storage_values"

	BidBid = "bid"
	BidLot = "lot"
	BidGuy = "guy"
	BidTic = "tic"
	BidEnd = "end"
)

var (
	BidsIndex = storage.IndexOne

	VatKey      = common.HexToHash(storage.IndexTwo)
	VatMetadata = utils.GetStorageValueMetadata(Vat, nil, utils.Address)

	GemKey      = common.HexToHash(storage.IndexThree)
	GemMetadata = utils.GetStorageValueMetadata(Gem, nil, utils.Address)

	BegKey      = common.HexToHash(storage.IndexFour)
	BegMetadata = utils.GetStorageValueMetadata(Beg, nil, utils.Uint256)

	TtlAndTauKey      = common.HexToHash(storage.IndexFive)
	packedTypes       = map[int]utils.ValueType{0: utils.Uint48, 1: utils.Uint48}
	packedNames       = map[int]string{0: Ttl, 1: Tau}
	TtlAndTauMetadata = utils.GetStorageValueMetadataForPackedSlot(Packed, nil, utils.PackedSlot, packedNames, packedTypes)

	KicksKey      = common.HexToHash(storage.IndexSix)
	KicksMetadata = utils.GetStorageValueMetadata(Kicks, nil, utils.Uint256)

	LiveKey      = common.HexToHash(storage.IndexSeven)
	LiveMetadata = utils.GetStorageValueMetadata(Live, nil, utils.Uint256)
)

type StorageKeysLookup struct {
	StorageRepository s2.IMakerStorageRepository
	ContractAddress   string
	mappings          map[common.Hash]utils.StorageValueMetadata
}

func (mappings StorageKeysLookup) Lookup(key common.Hash) (utils.StorageValueMetadata, error) {
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

func (mappings *StorageKeysLookup) SetDB(db *postgres.DB) {
	mappings.StorageRepository.SetDB(db)
}

func (mappings *StorageKeysLookup) loadMappings() error {
	mappings.mappings = loadStaticMappings()
	return mappings.loadBidKeys()
}

func loadStaticMappings() map[common.Hash]utils.StorageValueMetadata {
	mappings := make(map[common.Hash]utils.StorageValueMetadata)
	mappings[VatKey] = VatMetadata
	mappings[GemKey] = GemMetadata
	mappings[BegKey] = BegMetadata
	mappings[TtlAndTauKey] = TtlAndTauMetadata
	mappings[KicksKey] = KicksMetadata
	mappings[LiveKey] = LiveMetadata
	return mappings
}

func (mappings *StorageKeysLookup) loadBidKeys() error {
	bidIds, getBidsErr := mappings.StorageRepository.GetFlopBidIds(mappings.ContractAddress)
	if getBidsErr != nil {
		return getBidsErr
	}
	for _, bidId := range bidIds {
		hexBidId, convertErr := shared.ConvertIntStringToHex(bidId)
		if convertErr != nil {
			return convertErr
		}
		mappings.mappings[getBidBidKey(hexBidId)] = getBidBidMetadata(bidId)
		mappings.mappings[getBidLotKey(hexBidId)] = getBidLotMetadata(bidId)
		mappings.mappings[getBidGuyKey(hexBidId)] = getBidGuyMetadata(bidId)
		mappings.mappings[getBidTicKey(hexBidId)] = getBidTicMetadata(bidId)
		mappings.mappings[getBidEndKey(hexBidId)] = getBidEndMetadata(bidId)
	}
	return getBidsErr
}

func getBidBidKey(hexBidId string) common.Hash {
	return storage.GetMapping(BidsIndex, hexBidId)
}

func getBidBidMetadata(bidId string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.BidId: bidId}
	return utils.GetStorageValueMetadata(BidBid, keys, utils.Uint256)
}

func getBidLotKey(hexBidId string) common.Hash {
	return storage.GetIncrementedKey(getBidBidKey(hexBidId), 1)
}

func getBidLotMetadata(bidId string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.BidId: bidId}
	return utils.GetStorageValueMetadata(BidLot, keys, utils.Uint256)
}

func getBidGuyKey(hexBidId string) common.Hash {
	return storage.GetIncrementedKey(getBidBidKey(hexBidId), 2)
}

func getBidGuyMetadata(bidId string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.BidId: bidId}
	return utils.GetStorageValueMetadata(BidGuy, keys, utils.Address)
}

func getBidTicKey(hexBidId string) common.Hash {
	return storage.GetIncrementedKey(getBidBidKey(hexBidId), 3)
}

func getBidTicMetadata(bidId string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.BidId: bidId}
	return utils.GetStorageValueMetadata(BidTic, keys, utils.Uint48)
}

func getBidEndKey(hexBidId string) common.Hash {
	return storage.GetIncrementedKey(getBidBidKey(hexBidId), 4)
}

func getBidEndMetadata(bidId string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.BidId: bidId}
	return utils.GetStorageValueMetadata(BidEnd, keys, utils.Uint48)
}
