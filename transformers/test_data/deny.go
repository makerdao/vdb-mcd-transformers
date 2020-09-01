package test_data

import (
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
)

var rawDenyLog = types.Log{
	Address: common.HexToAddress(Cat100Address()),
	Topics: []common.Hash{
		common.HexToHash(constants.DenySignature()),
		common.HexToHash("0x00000000000000000000000013141b8a5e4a82ebc6b636849dd6a515185d6237"),
		common.HexToHash("0x00000000000000000000000013141b8a5e4a82ebc6b636849dd6a515185d6236"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e09c52a7f100000000000000000000000013141b8a5e4a82ebc6b636849dd6a515185d62360000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 10,
	TxHash:      common.HexToHash("0xe8f39fbb7fea3621f543868f19b1114e305aff6a063a30d32835ff1012526f91"),
	TxIndex:     7,
	BlockHash:   fakes.FakeHash,
	Index:       8,
	Removed:     false,
}

var DenyEventLog = core.EventLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawDenyLog,
	Transformed: false,
}

func DenyModel() event.InsertionModel { return CopyModel(denyModel) }

var denyModel = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.DenyTable,
	OrderedColumns: []event.ColumnName{
		event.HeaderFK,
		event.LogFK,
		event.AddressFK,
		constants.MsgSenderColumn,
		constants.UsrColumn,
	},
	ColumnValues: event.ColumnValues{event.HeaderFK: DenyEventLog.HeaderID, event.LogFK: DenyEventLog.ID},
}

var rawVatDenyLog = types.Log{
	Address: common.HexToAddress(VatAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.DenySignature()),
		common.HexToHash("0x00000000000000000000000013141b8a5e4a82ebc6b636849dd6a515185d6236"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e09c52a7f100000000000000000000000013141b8a5e4a82ebc6b636849dd6a515185d62360000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 10,
	TxHash:      common.HexToHash("0xe8f39fbb7fea3621f543868f19b1114e305aff6a063a30d32835ff1012526f91"),
	TxIndex:     7,
	BlockHash:   fakes.FakeHash,
	Index:       8,
	Removed:     false,
}

var VatDenyEventLog = core.EventLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawVatDenyLog,
	Transformed: false,
}

func VatDenyModel() event.InsertionModel { return CopyModel(vatDenyModel) }

var vatDenyModel = event.InsertionModel{
	SchemaName:     constants.MakerSchema,
	TableName:      constants.VatDenyTable,
	OrderedColumns: []event.ColumnName{event.HeaderFK, event.LogFK, constants.UsrColumn},
	ColumnValues: event.ColumnValues{
		event.HeaderFK: VatDenyEventLog.HeaderID,
		event.LogFK:    VatDenyEventLog.ID,
	},
}
