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

var rawPotDripLog = types.Log{
	Address: common.HexToAddress(PotAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.PotDripSignature()),
		common.HexToHash("0x00000000000000000000000087e76b0a50efc20259cafe0530f75ae0e816aaf2"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e09f678cca00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 15407281,
	TxHash:      common.HexToHash("0xe6077a8c1d42bc2cfe4f2b829cd66e75f9884677677615996f92a2b4d04110be"),
	TxIndex:     0,
	BlockHash:   common.HexToHash("0xf450da7f0a560375866bb0fc1d48fc59a43626afaae7e3e0183919dd6f412bb7"),
	Index:       1,
	Removed:     false,
}

var PotDripHeaderSyncLog = core.HeaderSyncLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawPotDripLog,
	Transformed: false,
}

var potDripModel = event.InsertionModel{
	SchemaName:     constants.MakerSchema,
	TableName:      constants.PotDripTable,
	OrderedColumns: []event.ColumnName{event.HeaderFK, event.LogFK, constants.MsgSenderColumn},
	ColumnValues: event.ColumnValues{
		event.HeaderFK: PotDripHeaderSyncLog.HeaderID,
		event.LogFK:    PotDripHeaderSyncLog.ID,
		// MsgSender column
	},
}

func PotDripModel() event.InsertionModel {
	return CopyModel(potDripModel)
}
