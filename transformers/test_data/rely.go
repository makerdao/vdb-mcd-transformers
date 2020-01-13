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

var rawRelyLog = types.Log{
	Address: common.HexToAddress(CatAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.RelySignature()),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x0000000000000000000000000e4725db88bb038bba4c4723e91ba183be11edf3"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e065fae35e0000000000000000000000000e4725db88bb038bba4c4723e91ba183be11edf30000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 10,
	TxHash:      common.HexToHash("0xe8f39fbb7fea3621f543868f19b1114e305aff6a063a30d32835ff1012526f91"),
	TxIndex:     7,
	BlockHash:   fakes.FakeHash,
	Index:       8,
	Removed:     false,
}

var RelyEventLog = core.EventLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawRelyLog,
	Transformed: false,
}

func RelyModel() event.InsertionModel { return CopyModel(relyModel) }

var relyModel = event.InsertionModel{
	SchemaName:     constants.MakerSchema,
	TableName:      constants.RelyTable,
	OrderedColumns: []event.ColumnName{event.HeaderFK, event.LogFK, event.AddressFK, constants.UsrColumn},
	ColumnValues:   event.ColumnValues{event.HeaderFK: RelyEventLog.HeaderID, event.LogFK: RelyEventLog.ID},
}

var rawVatRelyLog = types.Log{
	Address: common.HexToAddress(VatAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.RelySignature()),
		common.HexToHash("0x000000000000000000000000145b00b1ac4f01e84594efa2972fce1f5beb5ced"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e065fae35e000000000000000000000000145b00b1ac4f01e84594efa2972fce1f5beb5ced0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 10,
	TxHash:      common.HexToHash("0xe8f39fbb7fea3621f543868f19b1114e305aff6a063a30d32835ff1012526f91"),
	TxIndex:     7,
	BlockHash:   fakes.FakeHash,
	Index:       8,
	Removed:     false,
}

var VatRelyEventLog = core.EventLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawVatRelyLog,
	Transformed: false,
}

func VatRelyModel() event.InsertionModel { return CopyModel(vatRelyModel) }

var vatRelyModel = event.InsertionModel{
	SchemaName:     constants.MakerSchema,
	TableName:      constants.RelyTable,
	OrderedColumns: []event.ColumnName{event.HeaderFK, event.LogFK, event.AddressFK, constants.UsrColumn},
	ColumnValues: event.ColumnValues{
		event.HeaderFK: VatRelyEventLog.HeaderID,
		event.LogFK:    VatRelyEventLog.ID,
	},
}
