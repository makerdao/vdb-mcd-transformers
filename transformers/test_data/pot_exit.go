package test_data

import (
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
)

var rawPotExitLog = types.Log{
	Address: common.HexToAddress(PotAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.PotExitSignature()),
		common.HexToHash("0x00000000000000000000000071dd45d9579a499b58aa85f50e5e3b241ca2d10d"),
		common.HexToHash("0x00000000000000000000000000000000000000000000000ad2751323a735252d"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e07f8661a100000000000000000000000000000000000000000000000ad2751323a735252d0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 105,
	TxHash:      common.HexToHash("0xe5bebbe1ea46d8b6d1515ce9ac4659f9e6774669c1f2009dcc7289c18b91e393"),
	TxIndex:     2,
	BlockHash:   common.HexToHash("0x3e011e52723db56476dc8cd45e1325f7bf3f3b2d89651253d6e8b66489f37d7c"),
	Index:       3,
	Removed:     false,
}

var PotExitHeaderSyncLog = core.HeaderSyncLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawPotExitLog,
	Transformed: false,
}

func PotExitModel() event.InsertionModel { return CopyModel(potExitModel) }

var potExitModel = event.InsertionModel{
	SchemaName:     constants.MakerSchema,
	TableName:      constants.PotExitTable,
	OrderedColumns: []event.ColumnName{event.HeaderFK, event.LogFK, constants.MsgSenderColumn, constants.WadColumn},
	ColumnValues: event.ColumnValues{
		constants.WadColumn: "199632489101185590573",
		event.HeaderFK:      PotExitHeaderSyncLog.HeaderID,
		event.LogFK:         PotExitHeaderSyncLog.ID,
	},
}
