package test_data

import (
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/pot_file/dsr"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/pot_file/vow"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
)

var rawPotFileDSRLog = types.Log{
	Address: common.HexToAddress(PotAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.PotFileDSRSignature()),
		common.HexToHash("0x0000000000000000000000000e4725db88bb038bba4c4723e91ba183be11edf3"),
		common.HexToHash("0x6473720000000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x0000000000000000000000000000000000000000033b2e3ca88761c99baf1532"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e029ae811464737200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000033b2e3ca88761c99baf1532000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 14764699,
	TxHash:      common.HexToHash("0x6133d54cd52dbb29cd240631b963bf0db3976dbb83290b1b90829a371ab283f4"),
	TxIndex:     4,
	BlockHash:   common.HexToHash("0x75a0286bc101e6691b85fe1dc673beb46e08a4a662710846e677473cf0be54bd"),
	Index:       5,
	Removed:     false,
}

var rawPotFileVowLog = types.Log{
	Address: common.HexToAddress(PotAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.PotFileVowSignature()),
		common.HexToHash("0x00000000000000000000000013141b8a5e4a82ebc6b636849dd6a515185d6236"),
		common.HexToHash("0x766f770000000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x0000000000000000000000000f4cbe6cba918b7488c26e29d9ecd7368f38ea3b"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e0d4e8be83766f7700000000000000000000000000000000000000000000000000000000000000000000000000000000000f4cbe6cba918b7488c26e29d9ecd7368f38ea3b000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 14764543,
	TxHash:      common.HexToHash("0x77cb44b6811e4f81ee7d6a6cccdee8a525b0437c9dddb0e42257258749088c4d"),
	TxIndex:     1,
	BlockHash:   common.HexToHash("0x44bcf7b6bbbd71f329548928c03c6873a09c713bfb48c9aa086ae36a4eb3c611"),
	Index:       9,
	Removed:     false,
}

var PotFileDSRHeaderSyncLog = core.HeaderSyncLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawPotFileDSRLog,
	Transformed: false,
}

var PotFileVowHeaderSyncLog = core.HeaderSyncLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawPotFileVowLog,
	Transformed: false,
}

var potFileDSRModel = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.PotFileDSRTable,
	OrderedColumns: []event.ColumnName{
		constants.HeaderFK, constants.LogFK, dsr.What, dsr.Data,
	},
	ColumnValues: event.ColumnValues{
		constants.HeaderFK: PotFileDSRHeaderSyncLog.HeaderID,
		constants.LogFK:    PotFileDSRHeaderSyncLog.ID,
		dsr.What:           "dsr",
		dsr.Data:           "1000000000627937192491029810",
	},
}

var potFileVowModel = event.InsertionModel{
	SchemaName:     constants.MakerSchema,
	TableName:      constants.PotFileVowTable,
	OrderedColumns: []event.ColumnName{constants.HeaderFK, constants.LogFK, vow.What, vow.Data},
	ColumnValues: event.ColumnValues{
		constants.HeaderFK: PotFileVowHeaderSyncLog.HeaderID,
		constants.LogFK:    PotFileVowHeaderSyncLog.ID,
		vow.What:           "vow",
		vow.Data:           "0x0F4Cbe6CBA918b7488C26E29d9ECd7368F38EA3b",
	},
}

func PotFileDSRModel() event.InsertionModel { return CopyEventModel(potFileDSRModel) }
func PotFileVowModel() event.InsertionModel { return CopyEventModel(potFileVowModel) }
