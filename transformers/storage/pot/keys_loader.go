package pot

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities/wards"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	vdbStorage "github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

const (
	UserPie = "pie"
	Pie     = "Pie"
	Dsr     = "dsr"
	Chi     = "chi"
	Vat     = "vat"
	Vow     = "vow"
	Rho     = "rho"
	Live    = "live"
)

var (
	UserPieIndex = vdbStorage.IndexOne

	PieKey      = common.HexToHash(vdbStorage.IndexTwo)
	PieMetadata = types.GetValueMetadata(Pie, nil, types.Uint256)

	DsrKey      = common.HexToHash(vdbStorage.IndexThree)
	DsrMetadata = types.GetValueMetadata(Dsr, nil, types.Uint256)

	ChiKey      = common.HexToHash(vdbStorage.IndexFour)
	ChiMetadata = types.GetValueMetadata(Chi, nil, types.Uint256)

	VatKey      = common.HexToHash(vdbStorage.IndexFive)
	VatMetadata = types.GetValueMetadata(Vat, nil, types.Address)

	VowKey      = common.HexToHash(vdbStorage.IndexSix)
	VowMetadata = types.GetValueMetadata(Vow, nil, types.Address)

	RhoKey      = common.HexToHash(vdbStorage.IndexSeven)
	RhoMetadata = types.GetValueMetadata(Rho, nil, types.Uint256)

	LiveKey      = common.HexToHash(vdbStorage.IndexEight)
	LiveMetadata = types.GetValueMetadata(Live, nil, types.Uint256)
)

type keysLoader struct {
	storageRepository mcdStorage.IMakerStorageRepository
	contractAddress   string
}

func NewKeysLoader(storageRepository mcdStorage.IMakerStorageRepository, contractAddress string) storage.KeysLoader {
	return &keysLoader{storageRepository: storageRepository, contractAddress: contractAddress}
}

func (loader *keysLoader) SetDB(db *postgres.DB) {
	loader.storageRepository.SetDB(db)
}

func (loader *keysLoader) LoadMappings() (map[common.Hash]types.ValueMetadata, error) {
	mappings := getStaticMappings()
	users, err := loader.storageRepository.GetPotPieUsers()
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		paddedUser, padErr := utilities.PadAddress(user)
		if padErr != nil {
			return nil, padErr
		}
		mappings[getUserPieKey(paddedUser)] = getUserPieMetadata(user)
	}
	mappings, wardsErr := loader.addWardsKeys(mappings)
	if wardsErr != nil {
		return nil, wardsErr
	}
	return mappings, nil
}

func (loader *keysLoader) addWardsKeys(mappings map[common.Hash]types.ValueMetadata) (map[common.Hash]types.ValueMetadata, error) {
	addresses, err := loader.storageRepository.GetWardsAddresses(loader.contractAddress)
	if err != nil {
		return nil, err
	}
	return wards.AddWardsKeys(mappings, addresses)
}

func getStaticMappings() map[common.Hash]types.ValueMetadata {
	mappings := make(map[common.Hash]types.ValueMetadata)
	mappings[PieKey] = PieMetadata
	mappings[DsrKey] = DsrMetadata
	mappings[ChiKey] = ChiMetadata
	mappings[VatKey] = VatMetadata
	mappings[VowKey] = VowMetadata
	mappings[RhoKey] = RhoMetadata
	mappings[LiveKey] = LiveMetadata
	return mappings
}

func getUserPieKey(user string) common.Hash {
	return vdbStorage.GetKeyForMapping(UserPieIndex, user)
}

func getUserPieMetadata(user string) types.ValueMetadata {
	keys := map[types.Key]string{constants.MsgSender: user}
	return types.GetValueMetadata(UserPie, keys, types.Uint256)
}
