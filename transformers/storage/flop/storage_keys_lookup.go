package flop

import (
	"github.com/ethereum/go-ethereum/common"
	vdbStorage "github.com/vulcanize/vulcanizedb/libraries/shared/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage"
)

var (
	BidsIndex = vdbStorage.IndexOne

	VatKey      = common.HexToHash(vdbStorage.IndexTwo)
	VatMetadata = utils.GetStorageValueMetadata(storage.Vat, nil, utils.Address)

	GemKey      = common.HexToHash(vdbStorage.IndexThree)
	GemMetadata = utils.GetStorageValueMetadata(storage.Gem, nil, utils.Address)

	BegKey      = common.HexToHash(vdbStorage.IndexFour)
	BegMetadata = utils.GetStorageValueMetadata(storage.Beg, nil, utils.Uint256)

	TtlAndTauKey      = common.HexToHash(vdbStorage.IndexFive)
	ttlAndTauTypes    = map[int]utils.ValueType{0: utils.Uint48, 1: utils.Uint48}
	ttlAndTauNames    = map[int]string{0: storage.Ttl, 1: storage.Tau}
	TtlAndTauMetadata = utils.GetStorageValueMetadataForPackedSlot(storage.Packed, nil, utils.PackedSlot, ttlAndTauNames, ttlAndTauTypes)

	KicksKey      = common.HexToHash(vdbStorage.IndexSix)
	KicksMetadata = utils.GetStorageValueMetadata(storage.Kicks, nil, utils.Uint256)

	LiveKey      = common.HexToHash(vdbStorage.IndexSeven)
	LiveMetadata = utils.GetStorageValueMetadata(storage.Live, nil, utils.Uint256)
)

type StorageKeysLookup struct {
	StorageRepository storage.IMakerStorageRepository
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
	err := mappings.loadBidKeys()
	if err != nil {
		return err
	}
	mappings.mappings = vdbStorage.AddHashedKeys(mappings.mappings)
	return nil
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
		mappings.mappings[getBidGuyTicEndKey(hexBidId)] = getBidGuyTicEndMetadata(bidId)
	}
	return getBidsErr
}

func getBidBidKey(hexBidId string) common.Hash {
	return vdbStorage.GetMapping(BidsIndex, hexBidId)
}

func getBidBidMetadata(bidId string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.BidId: bidId}
	return utils.GetStorageValueMetadata(storage.BidBid, keys, utils.Uint256)
}

func getBidLotKey(hexBidId string) common.Hash {
	return vdbStorage.GetIncrementedKey(getBidBidKey(hexBidId), 1)
}

func getBidLotMetadata(bidId string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.BidId: bidId}
	return utils.GetStorageValueMetadata(storage.BidLot, keys, utils.Uint256)
}

func getBidGuyTicEndKey(hexBidId string) common.Hash {
	return vdbStorage.GetIncrementedKey(getBidBidKey(hexBidId), 2)
}

func getBidGuyTicEndMetadata(bidId string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.BidId: bidId}
	packedTypes := map[int]utils.ValueType{0: utils.Address, 1: utils.Uint48, 2: utils.Uint48}
	packedNames := map[int]string{0: storage.BidGuy, 1: storage.BidTic, 2: storage.BidEnd}
	return utils.GetStorageValueMetadataForPackedSlot(storage.Packed, keys, utils.PackedSlot, packedNames, packedTypes)
}
