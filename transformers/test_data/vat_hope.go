package test_data

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"math/rand"
)

var rawVatHopeLog = types.Log{
	Address: common.HexToAddress(VatAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.VatHopeSignature()),
		common.HexToHash("0x0000000000000000000000005ef30b9986345249bc32d8928b7ee64de9435e39"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e0a3b22fc40000000000000000000000005ef30b9986345249bc32d8928b7ee64de9435e390000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 9861938,
	TxHash:      common.HexToHash("0x76931cb92fa31ebf5451daed262d01736b1ac3f4076f4e8428e6cdfe940a8129"),
	TxIndex:     317,
	BlockHash:   common.HexToHash("0x0843d32f14e1b258a511c187270018b957bd566787afb396c78e1455305322d0"),
	Index:       61,
	Removed:     false,
}

var VatHopeEventLog = core.EventLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawVatHopeLog,
	Transformed: false,
}

var vatHopeModel = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.VatHopeTable,
	OrderedColumns: []event.ColumnName{
		event.HeaderFK, event.LogFK, constants.UsrColumn,
	},
	ColumnValues: event.ColumnValues{
		event.HeaderFK: VatHopeEventLog.HeaderID,
		event.LogFK:    VatHopeEventLog.ID,
		//constants.UsrColumn
	},
}

func VatHopeModel() event.InsertionModel { return CopyModel(vatHopeModel) }
