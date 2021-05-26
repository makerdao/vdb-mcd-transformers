package clip

import (
	"github.com/ethereum/go-ethereum/common"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	vdbStorage "github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

var (
	//wards at index0
	IlkKey      = common.HexToHash(vdbStorage.IndexOne)
	IlkMetadata = types.GetValueMetadata(mcdStorage.Ilk, nil, types.Bytes32)

	VatKey      = common.HexToHash(vdbStorage.IndexTwo)
	VatMetadata = types.GetValueMetadata(mcdStorage.Vat, nil, types.Address)

	DogKey      = common.HexToHash(vdbStorage.IndexThree)
	DogMetadata = types.GetValueMetadata(mcdStorage.Dog, nil, types.Address)

	VowKey      = common.HexToHash(vdbStorage.IndexFour)
	VowMetadata = types.GetValueMetadata(mcdStorage.Vow, nil, types.Address)

	SpotterKey      = common.HexToHash(vdbStorage.IndexFive)
	SpotterMetadata = types.GetValueMetadata(mcdStorage.Spotter, nil, types.Address)

	CalcKey      = common.HexToHash(vdbStorage.IndexSix)
	CalcMetadata = types.GetValueMetadata(mcdStorage.Calc, nil, types.Address)

	BufKey      = common.HexToHash(vdbStorage.IndexSeven)
	BufMetadata = types.GetValueMetadata(mcdStorage.Buf, nil, types.Uint256)

	TailKey      = common.HexToHash(vdbStorage.IndexEight)
	TailMetadata = types.GetValueMetadata(mcdStorage.Tail, nil, types.Uint256)

	CuspKey      = common.HexToHash(vdbStorage.IndexNine)
	CuspMetadata = types.GetValueMetadata(mcdStorage.Cusp, nil, types.Uint256)

	// TODO: Add actual types to vulcanizedb (uint64 and uint192)

	ChipAndTipStorageKey = common.HexToHash(vdbStorage.IndexTen)
	chipAndTipTypes      = map[int]types.ValueType{0: types.Uint32, 1: types.Uint48}
	chipAndTipNames      = map[int]string{0: mcdStorage.Chip, 1: mcdStorage.Tip}
	ChipAndTipMetadata   = types.GetValueMetadataForPackedSlot(mcdStorage.Packed, nil, types.PackedSlot, chipAndTipNames, chipAndTipTypes)

	ChostKey      = common.HexToHash(vdbStorage.IndexEleven)
	ChostMetadata = types.GetValueMetadata(mcdStorage.Chost, nil, types.Uint256)

	KicksKey      = common.HexToHash(vdbStorage.IndexTwelve)
	KicksMetadata = types.GetValueMetadata(mcdStorage.Kicks, nil, types.Uint256)

	ActiveKey	  = common.HexToHash(vdbStorage.IndexThirteen)
	ActiveMetadata = types.GetValueMetadata(mcdStorage.Active, nil, types.Uint256)
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
	return mappings, nil
}

func loadStaticMappings() map[common.Hash]types.ValueMetadata {
	mappings := make(map[common.Hash]types.ValueMetadata)
	mappings[VatKey] = VatMetadata
	mappings[IlkKey] = IlkMetadata
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
