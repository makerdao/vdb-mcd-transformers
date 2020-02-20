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

var rawOsmChangeLog = types.Log{
	Address: common.HexToAddress(EthOsmAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.OsmChangeSignature()),
		common.HexToHash("0x000000000000000000000000dc984d513a0f9ca9aa602d4df8517677918936e3"),
		common.HexToHash("0x000000000000000000000000a950524441892a31ebddf91d3ceefa04bf454466"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e03b663195434f4c352d4100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 10980181,
	TxHash:      common.HexToHash("0x27f7834f778ec7d4289cf3337f8e428785c6d023164c02fc44565dbf2e26c49a"),
	TxIndex:     10,
	BlockHash:   fakes.FakeHash,
	Index:       11,
	Removed:     false,
}

var OsmChangeEventLog = core.EventLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawOsmChangeLog,
	Transformed: false,
}

func OsmChangeModel() event.InsertionModel { return CopyModel(osmChangeModel) }

var osmChangeModel = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.OsmChangeTable,
	OrderedColumns: []event.ColumnName{
		event.HeaderFK,
		event.LogFK,
		event.AddressFK,
		constants.MsgSenderColumn,
		constants.SrcColumn,
	},
	ColumnValues: event.ColumnValues{
		event.HeaderFK: OsmChangeEventLog.HeaderID,
		event.LogFK:    OsmChangeEventLog.ID,
	},
}
