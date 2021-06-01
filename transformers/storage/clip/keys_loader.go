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

	ActiveKey      = common.HexToHash(vdbStorage.IndexEleven)
	ActiveMetadata = types.GetValueMetadata(Active, nil, types.Uint256)
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
	mappings[ActiveKey] = ActiveMetadata
	return mappings
}
