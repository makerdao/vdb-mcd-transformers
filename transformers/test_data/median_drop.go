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

var rawMedianDropLog = types.Log{
	Address: common.HexToAddress(EthMedianAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.MedianDropSignature()),
		common.HexToHash("0x000000000000000000000000c45e7858eef1318337a803ede8c5a9be12e2b40f"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000020"),
		common.HexToHash("0x000000000000000000000000000000000000000000000000000000000000000a"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e08ef5eaf00000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000c45e7858eef1318337a803ede8c5a9be12e2b40f000000000000000000000000ef6b95815e215635bd77851f1fc42e87508730240000000000000000000000008efccc4ecb27f7f233a7ff4e74e86c5e979d1c43000000000000000000000000c2d2d553a39cc08e7e294427ede2c38a89c0066a00000000000000000000000036c7d1aee129f32a07609a03dc5ffff6"),
	BlockNumber: 8936530,
	TxHash:      common.HexToHash("0x27f7834f778ec7d4289cf3337f8e428785c6d023164c02fc44565dbf2e26c49a"),
	TxIndex:     10,
	BlockHash:   fakes.FakeHash,
	Index:       11,
	Removed:     false,
}

var MedianDropLog = core.EventLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawMedianDropLog,
	Transformed: false,
}

func MedianDropModel() event.InsertionModel { return CopyModel(medianDropModel) }

var medianDropModel = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.MedianDropTable,
	OrderedColumns: []event.ColumnName{
		event.HeaderFK,
		event.LogFK,
		event.AddressFK,
		constants.MsgSenderColumn,
		constants.AColumn,
		constants.A2Column,
		constants.A3Column,
		constants.A4Column,
	},
	ColumnValues: event.ColumnValues{
		event.HeaderFK: MedianDropLog.HeaderID,
		event.LogFK:    MedianDropLog.ID,
	},
}
