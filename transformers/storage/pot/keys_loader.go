package pot

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	vdbStorage "github.com/makerdao/vulcanizedb/libraries/shared/storage"
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
	PieMetadata = vdbStorage.GetValueMetadata(Pie, nil, vdbStorage.Uint256)

	DsrKey      = common.HexToHash(vdbStorage.IndexThree)
	DsrMetadata = vdbStorage.GetValueMetadata(Dsr, nil, vdbStorage.Uint256)

	ChiKey      = common.HexToHash(vdbStorage.IndexFour)
	ChiMetadata = vdbStorage.GetValueMetadata(Chi, nil, vdbStorage.Uint256)

	VatKey      = common.HexToHash(vdbStorage.IndexFive)
	VatMetadata = vdbStorage.GetValueMetadata(Vat, nil, vdbStorage.Address)

	VowKey      = common.HexToHash(vdbStorage.IndexSix)
	VowMetadata = vdbStorage.GetValueMetadata(Vow, nil, vdbStorage.Address)

	RhoKey      = common.HexToHash(vdbStorage.IndexSeven)
	RhoMetadata = vdbStorage.GetValueMetadata(Rho, nil, vdbStorage.Uint256)

	LiveKey      = common.HexToHash(vdbStorage.IndexEight)
	LiveMetadata = vdbStorage.GetValueMetadata(Live, nil, vdbStorage.Uint256)
)

type keysLoader struct {
	storageRepository mcdStorage.IMakerStorageRepository
}

func NewKeysLoader(storageRepository mcdStorage.IMakerStorageRepository) storage.KeysLoader {
	return &keysLoader{storageRepository: storageRepository}
}

func (loader *keysLoader) SetDB(db *postgres.DB) {
	loader.storageRepository.SetDB(db)
}

func (loader *keysLoader) LoadMappings() (map[common.Hash]vdbStorage.ValueMetadata, error) {
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
	return mappings, nil
}

func getStaticMappings() map[common.Hash]vdbStorage.ValueMetadata {
	mappings := make(map[common.Hash]vdbStorage.ValueMetadata)
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

func getUserPieMetadata(user string) vdbStorage.ValueMetadata {
	keys := map[vdbStorage.Key]string{constants.MsgSender: user}
	return vdbStorage.GetValueMetadata(UserPie, keys, vdbStorage.Uint256)
}
