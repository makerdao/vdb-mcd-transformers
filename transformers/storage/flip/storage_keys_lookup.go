package flip

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
	BidsMappingIndex = vdbStorage.IndexOne

	VatKey      = common.HexToHash(vdbStorage.IndexTwo)
	VatMetadata = utils.GetStorageValueMetadata(storage.Vat, nil, utils.Address)

	IlkKey      = common.HexToHash(vdbStorage.IndexThree)
	IlkMetadata = utils.GetStorageValueMetadata(storage.Ilk, nil, utils.Bytes32)

	BegKey      = common.HexToHash(vdbStorage.IndexFour)
	BegMetadata = utils.GetStorageValueMetadata(storage.Beg, nil, utils.Uint256)

	TtlAndTauStorageKey = common.HexToHash(vdbStorage.IndexFive)
	ttlAndTauTypes      = map[int]utils.ValueType{0: utils.Uint48, 1: utils.Uint48}
	ttlAndTauNames      = map[int]string{0: storage.Ttl, 1: storage.Tau}
	TtlAndTauMetadata   = utils.GetStorageValueMetadataForPackedSlot(storage.Packed, nil, utils.PackedSlot, ttlAndTauNames, ttlAndTauTypes)

	KicksKey      = common.HexToHash(vdbStorage.IndexSix)
	KicksMetadata = utils.GetStorageValueMetadata(storage.Kicks, nil, utils.Uint256)
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
	mappings[IlkKey] = IlkMetadata
	mappings[BegKey] = BegMetadata
	mappings[TtlAndTauStorageKey] = TtlAndTauMetadata
	mappings[KicksKey] = KicksMetadata
	return mappings
}

func (mappings *StorageKeysLookup) loadBidKeys() error {
	bidIds, err := mappings.StorageRepository.GetFlipBidIds(mappings.ContractAddress)
	if err != nil {
		return err
	}
	for _, bidId := range bidIds {
		hexBidId, err := shared.ConvertIntStringToHex(bidId)
		if err != nil {
			return err
		}
		mappings.mappings[getBidBidKey(hexBidId)] = getBidBidMetadata(bidId)
		mappings.mappings[getBidLotKey(hexBidId)] = getBidLotMetadata(bidId)
		mappings.mappings[getBidGuyTicEndKey(hexBidId)] = getBidGuyTicEndMetadata(bidId)
		mappings.mappings[getBidUsrKey(hexBidId)] = getBidUsrMetadata(bidId)
		mappings.mappings[getBidGalKey(hexBidId)] = getBidGalMetadata(bidId)
		mappings.mappings[getBidTabKey(hexBidId)] = getBidTabMetadata(bidId)
	}
	return nil
}

func getBidBidKey(hexBidId string) common.Hash {
	return vdbStorage.GetMapping(BidsMappingIndex, hexBidId)
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

func getBidUsrKey(hexBidId string) common.Hash {
	return vdbStorage.GetIncrementedKey(getBidBidKey(hexBidId), 3)
}

func getBidUsrMetadata(bidId string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.BidId: bidId}
	return utils.GetStorageValueMetadata(storage.BidUsr, keys, utils.Address)
}

func getBidGalKey(hexBidId string) common.Hash {
	return vdbStorage.GetIncrementedKey(getBidBidKey(hexBidId), 4)
}

func getBidGalMetadata(bidId string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.BidId: bidId}
	return utils.GetStorageValueMetadata(storage.BidGal, keys, utils.Address)
}

func getBidTabKey(hexBidId string) common.Hash {
	return vdbStorage.GetIncrementedKey(getBidBidKey(hexBidId), 5)
}

func getBidTabMetadata(bidId string) utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.BidId: bidId}
	return utils.GetStorageValueMetadata(storage.BidTab, keys, utils.Uint256)
}
