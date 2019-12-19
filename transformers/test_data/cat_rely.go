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

var rawCatRelyLog = types.Log{
	Address: common.HexToAddress(CatAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.CatRelySignature()),
		common.HexToHash("0x00000000000000000000000039ad5d336a4c08fac74879f796e1ea0af26c1521"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e06924500900000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 14764546,
	TxHash:      common.HexToHash("0x16a1e5402ce5e4006f126ac4dd1d67a2a1e366c681fa8eec0e431041623b1fca"),
	TxIndex:     4,
	BlockHash:   common.HexToHash("0x774d2ae5736c5936775671e1945d0dde88900417e0c76ce3e50f23d434d6cff6"),
	Index:       5,
	Removed:     false,
}

var CatRelyHeaderSyncLog = core.HeaderSyncLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawCatRelyLog,
	Transformed: false,
}

var catRelyModel = event.InsertionModel{
	SchemaName:     constants.MakerSchema,
	TableName:      constants.CatRelyTable,
	OrderedColumns: []event.ColumnName{event.HeaderFK, event.LogFK, constants.UsrColumn},
	ColumnValues: event.ColumnValues{
		event.HeaderFK:      CatRelyHeaderSyncLog.HeaderID,
		event.LogFK:         CatRelyHeaderSyncLog.ID,
		constants.UsrColumn: "329278621794589981598016975744453791986620306721",
	},
}

func CatRelyModel() event.InsertionModel { return CopyModel(catRelyModel) }
