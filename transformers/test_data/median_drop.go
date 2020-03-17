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
		common.HexToHash("0x000000000000000000000000ddb108893104de4e1c6d0e47c42237db4e617acc"),
		common.HexToHash("0x000000000000000000000000b4eb54af9cc7882df0121d26c5b97e802915abe6"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e03b663195434f4c352d4100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
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
	},
	ColumnValues: event.ColumnValues{
		event.HeaderFK: MedianDropLog.HeaderID,
		event.LogFK:    MedianDropLog.ID,
	},
}
