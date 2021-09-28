package clip

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities/wards"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	vdbStorage "github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

var (

	//wards at index0
	DogKey      = common.HexToHash(vdbStorage.IndexOne)
	DogMetadata = types.GetValueMetadata(Dog, nil, types.Address)

	VowKey      = common.HexToHash(vdbStorage.IndexTwo)
	VowMetadata = types.GetValueMetadata(mcdStorage.Vow, nil, types.Address)

	SpotterKey      = common.HexToHash(vdbStorage.IndexThree)
	SpotterMetadata = types.GetValueMetadata(Spotter, nil, types.Address)

	CalcKey      = common.HexToHash(vdbStorage.IndexFour)
	CalcMetadata = types.GetValueMetadata(Calc, nil, types.Address)

	BufKey      = common.HexToHash(vdbStorage.IndexFive)
	BufMetadata = types.GetValueMetadata(Buf, nil, types.Uint256)

	TailKey      = common.HexToHash(vdbStorage.IndexSix)
	TailMetadata = types.GetValueMetadata(Tail, nil, types.Uint256)

	CuspKey      = common.HexToHash(vdbStorage.IndexSeven)
	CuspMetadata = types.GetValueMetadata(Cusp, nil, types.Uint256)

	ChipAndTipStorageKey = common.HexToHash(vdbStorage.IndexEight)
	chipAndTipTypes      = map[int]types.ValueType{0: types.Uint64, 1: types.Uint192}
	chipAndTipNames      = map[int]string{0: Chip, 1: Tip}
	ChipAndTipMetadata   = types.GetValueMetadataForPackedSlot(mcdStorage.Packed, nil, types.PackedSlot, chipAndTipNames, chipAndTipTypes)

	ChostKey      = common.HexToHash(vdbStorage.IndexNine)
	ChostMetadata = types.GetValueMetadata(Chost, nil, types.Uint256)

	KicksKey      = common.HexToHash(vdbStorage.IndexTen)
	KicksMetadata = types.GetValueMetadata(mcdStorage.Kicks, nil, types.Uint256)

	SalesMappingIndex = vdbStorage.IndexTwelve
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
		return nil, fmt.Errorf("error adding wards keys to clip keys loader: %w", wardsErr)
	}
	mappings, salesErr := loader.loadSalesKeys(mappings)
	if salesErr != nil {
		return nil, fmt.Errorf("error adding sales keys to clip keys loader: %w", salesErr)
	}
	return mappings, nil
}

func loadStaticMappings() map[common.Hash]types.ValueMetadata {
	mappings := make(map[common.Hash]types.ValueMetadata)
	mappings[DogKey] = DogMetadata
	mappings[VowKey] = VowMetadata
	mappings[SpotterKey] = SpotterMetadata
	mappings[CalcKey] = CalcMetadata
	mappings[BufKey] = BufMetadata
	mappings[TailKey] = TailMetadata
	mappings[CuspKey] = CuspMetadata
	mappings[ChipAndTipStorageKey] = ChipAndTipMetadata
	mappings[ChostKey] = ChostMetadata
	mappings[KicksKey] = KicksMetadata
	return mappings
}

func (loader *keysLoader) loadSalesKeys(mappings map[common.Hash]types.ValueMetadata) (map[common.Hash]types.ValueMetadata, error) {
	saleIDs, saleErr := loader.storageRepository.GetClipSalesIDs(loader.contractAddress)
	if saleErr != nil {
		return nil, fmt.Errorf("error getting clip sale IDs: %w", saleErr)
	}
	for _, saleID := range saleIDs {
		hexSaleID, convertErr := shared.ConvertIntStringToHex(saleID)
		if convertErr != nil {
			return nil, fmt.Errorf("error converting int string to hex: %w", convertErr)
		}
		mappings[getSalePosKey(hexSaleID)] = getSalePosMetadata(saleID)
		mappings[getSaleTabKey(hexSaleID)] = getSaleTabMetadata(saleID)
		mappings[getSaleLotKey(hexSaleID)] = getSaleLotMetadata(saleID)
		mappings[getSaleUsrTicKey(hexSaleID)] = getSaleUsrTicMetadata(saleID)
		mappings[getSaleTopKey(hexSaleID)] = getSaleTopMetadata(saleID)
	}
	return mappings, nil
}

func getSalePosKey(hexSaleID string) common.Hash {
	return vdbStorage.GetKeyForMapping(SalesMappingIndex, hexSaleID)
}

func getSalePosMetadata(saleID string) types.ValueMetadata {
	keys := map[types.Key]string{constants.SaleId: saleID}
	return types.GetValueMetadata(SalePos, keys, types.Uint256)
}

func getSaleTabKey(hexSaleID string) common.Hash {
	return vdbStorage.GetIncrementedKey(getSalePosKey(hexSaleID), 1)
}

func getSaleTabMetadata(saleID string) types.ValueMetadata {
	keys := map[types.Key]string{constants.SaleId: saleID}
	return types.GetValueMetadata(SaleTab, keys, types.Uint256)
}

func getSaleLotKey(hexSaleID string) common.Hash {
	return vdbStorage.GetIncrementedKey(getSalePosKey(hexSaleID), 2)
}

func getSaleLotMetadata(saleID string) types.ValueMetadata {
	keys := map[types.Key]string{constants.SaleId: saleID}
	return types.GetValueMetadata(SaleLot, keys, types.Uint256)
}

func getSaleUsrTicKey(hexSaleID string) common.Hash {
	return vdbStorage.GetIncrementedKey(getSalePosKey(hexSaleID), 3)
}

func getSaleUsrTicMetadata(saleID string) types.ValueMetadata {
	keys := map[types.Key]string{constants.SaleId: saleID}
	packedTypes := map[int]types.ValueType{0: types.Address, 1: types.Uint96}
	packedNames := map[int]string{0: SaleUsr, 1: SaleTic}
	return types.GetValueMetadataForPackedSlot(mcdStorage.Packed, keys, types.PackedSlot, packedNames, packedTypes)
}

func getSaleTopKey(hexSaleID string) common.Hash {
	return vdbStorage.GetIncrementedKey(getSalePosKey(hexSaleID), 4)
}

func getSaleTopMetadata(saleID string) types.ValueMetadata {
	keys := map[types.Key]string{constants.SaleId: saleID}
	return types.GetValueMetadata(SaleTop, keys, types.Uint256)
}

func addWardsKeys(mappings map[common.Hash]types.ValueMetadata, address string, repository mcdStorage.IMakerStorageRepository) (map[common.Hash]types.ValueMetadata, error) {
	addresses, err := repository.GetWardsAddresses(address)
	if err != nil {
		return nil, fmt.Errorf("error getting wards addresses: %w", err)
	}
	return wards.AddWardsKeys(mappings, addresses)
}
