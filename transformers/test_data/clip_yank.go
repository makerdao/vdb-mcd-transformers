package test_data

import (
	"math/big"
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
)

var saleID = big.NewInt(456)

var RawClipYankLog = types.Log{
	Address: common.HexToAddress(Clip1xxAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.ClipYankSignature()),
	},
	Data: hexutil.MustDecode("0x" +
		"00000000000000000000000000000000000000000000000000000000000001c8"), //id 456
	BlockNumber: uint64(testBlockNumber),
	TxHash:      common.HexToHash(flipKickTransactionHash),
	TxIndex:     999,
	BlockHash:   fakes.FakeHash,
	Index:       1,
	Removed:     false,
}

var ClipYankEventLog = core.EventLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         RawClipYankLog,
	Transformed: false,
}

var clipYankModel = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.ClipYankTable,
	OrderedColumns: []event.ColumnName{
		event.HeaderFK,
		event.LogFK,
		event.AddressFK,
		constants.SaleIDColumn,
	},
	ColumnValues: event.ColumnValues{
		event.HeaderFK: ClipYankEventLog.HeaderID,
		event.LogFK:    ClipYankEventLog.ID,
		// event.AddressFK
		constants.SaleIDColumn: saleID.String(),
	},
}

func ClipYankModel() event.InsertionModel { return CopyModel(clipYankModel) }
