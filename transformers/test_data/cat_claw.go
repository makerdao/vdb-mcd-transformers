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

var (
	rawCatClawLog = types.Log{
		Address: common.HexToAddress(Cat110Address()),
		Topics: []common.Hash{
			common.HexToHash("0xe66d279b00000000000000000000000000000000000000000000000000000000"),
			common.HexToHash("0x000000000000000000000000f32836b9e1f47a0515c6ec431592d5ebc276407f"),
			common.HexToHash("0x0000000000000000000000001ce1660f8ed632ec18845c77938bfd09a16743dc"),
			common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
		},
		Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e0e66d279b0000000000000000000000001ce1660f8ed632ec18845c77938bfd09a16743dc0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
		BlockNumber: 10773034,
		TxHash:      common.HexToHash("0xc1fcf14936867b312f87e28879e1e2ac0fd35ecd156004ed204c8b55e183c1ce"),
		TxIndex:     0,
		BlockHash:   common.HexToHash("0xd50db3d21a6ce21b66949136019209c364f81183ff280138934a846b96d86a7e"),
		Index:       0,
		Removed:     false,
	}

	CatClawEventLog = core.EventLog{
		ID:          int64(rand.Int31()),
		HeaderID:    int64(rand.Int31()),
		Log:         rawCatClawLog,
		Transformed: false,
	}

	catClawModel = event.InsertionModel{
		SchemaName: constants.MakerSchema,
		TableName:  constants.CatClawTable,
		OrderedColumns: []event.ColumnName{
			event.HeaderFK,
			event.AddressFK,
			constants.MsgSenderColumn,
			event.LogFK,
			constants.RadColumn,
		},
		ColumnValues: event.ColumnValues{
			event.HeaderFK:  CatClawEventLog.HeaderID,
			event.AddressFK: CatClawEventLog.Log.Address.Hex(),
			//constants.MsgSenderColumn
			event.LogFK:         CatClawEventLog.ID,
			constants.RadColumn: "164878299999999999999707170900220955677367419868",
		},
	}
)

func CatClawModel() event.InsertionModel { return CopyModel(catClawModel) }
