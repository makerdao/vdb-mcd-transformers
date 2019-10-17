package flap

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

	VatStorageKey = common.HexToHash(vdbStorage.IndexTwo)
	VatMetadata   = utils.GetStorageValueMetadata(storage.Vat, nil, utils.Address)

	GemStorageKey = common.HexToHash(vdbStorage.IndexThree)
	GemMetadata   = utils.GetStorageValueMetadata(storage.Gem, nil, utils.Address)

	BegStorageKey = common.HexToHash(vdbStorage.IndexFour)
	BegMetadata   = utils.GetStorageValueMetadata(storage.Beg, nil, utils.Uint256)

	TtlAndTauStorageKey = common.HexToHash(vdbStorage.IndexFive)
	ttlAndTauTypes      = map[int]utils.ValueType{0: utils.Uint48, 1: utils.Uint48}
	ttlAndTauNames      = map[int]string{0: storage.Ttl, 1: storage.Tau}
	TtlAndTauMetadata   = utils.GetStorageValueMetadataForPackedSlot(storage.Packed, nil, utils.PackedSlot, ttlAndTauNames, ttlAndTauTypes)

	KicksStorageKey = common.HexToHash(vdbStorage.IndexSix)
	KicksMetadata   = utils.GetStorageValueMetadata(storage.Kicks, nil, utils.Uint256)

	LiveStorageKey = common.HexToHash(vdbStorage.IndexSeven)
	LiveMetadata   = utils.GetStorageValueMetadata(storage.Live, nil, utils.Uint256)
)

type StorageKeysLookup struct {
	StorageRepository storage.IMakerStorageRepository
	mappings          map[common.Hash]utils.StorageValueMetadata
	ContractAddress   string
}

func (lookup *StorageKeysLookup) Lookup(key common.Hash) (utils.StorageValueMetadata, error) {
	metadata, ok := lookup.mappings[key]
	if !ok {
		loadErr := lookup.loadMapping()
		if loadErr != nil {
			return utils.StorageValueMetadata{}, loadErr
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

func (lookup *StorageKeysLookup) loadMapping() error {
	lookup.loadStaticKeys()
	err := lookup.loadBidKeys()
	if err != nil {
		return err
	}
	lookup.mappings = vdbStorage.AddHashedKeys(lookup.mappings)
	return nil
}

func (lookup *StorageKeysLookup) loadBidKeys() error {
	bidIds, getBidIdsErr := lookup.StorageRepository.GetFlapBidIds(lookup.ContractAddress)
	for _, bidId := range bidIds {
		hexBidId, convertErr := shared.ConvertIntStringToHex(bidId)
		if convertErr != nil {
			return convertErr
		}

		lookup.mappings[getBidBidKey(hexBidId)] = getBidBidMetadata(bidId)
		lookup.mappings[getBidLotKey(hexBidId)] = getBidLotMetadata(bidId)
		lookup.mappings[getBidGuyTicEndKey(hexBidId)] = getBidGuyTicEndMetadata(bidId)
	}

	return getBidIdsErr
}

func getBidBidKey(bidId string) common.Hash {
	return vdbStorage.GetMapping(BidsIndex, bidId)
}

func getBidBidMetadata(bidId string) utils.StorageValueMetadata {
	return utils.StorageValueMetadata{
		Name: storage.BidBid,
		Keys: map[utils.Key]string{constants.BidId: bidId},
		Type: utils.Uint256,
	}
}

func getBidLotKey(bidId string) common.Hash {
	return vdbStorage.GetIncrementedKey(getBidBidKey(bidId), 1) //should this be renamed GetMappingKey?
}

func getBidLotMetadata(bidId string) utils.StorageValueMetadata {
	return utils.StorageValueMetadata{
		Name: storage.BidLot,
		Keys: map[utils.Key]string{constants.BidId: bidId},
		Type: utils.Uint256,
	}
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

func (lookup *StorageKeysLookup) loadStaticKeys() {
	mappings := make(map[common.Hash]utils.StorageValueMetadata)
	mappings[VatStorageKey] = VatMetadata
	mappings[GemStorageKey] = GemMetadata
	mappings[BegStorageKey] = BegMetadata
	mappings[TtlAndTauStorageKey] = TtlAndTauMetadata
	mappings[KicksStorageKey] = KicksMetadata
	mappings[LiveStorageKey] = LiveMetadata
	lookup.mappings = mappings
}
