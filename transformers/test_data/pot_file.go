package test_data

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/pot_file/dsr"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"math/rand"
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

var PotFileDSRHeaderSyncLog = core.HeaderSyncLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawPotFileDSRLog,
	Transformed: false,
}

var potFileDSRModel = event.InsertionModel{
	SchemaName: "maker",
	TableName:  "pot_file_dsr",
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

func PotFileDSRModel() event.InsertionModel { return CopyEventModel(potFileDSRModel) }
