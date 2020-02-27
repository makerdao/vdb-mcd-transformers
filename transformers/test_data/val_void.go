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

var rawValVoidLog = types.Log{
	Address: common.HexToAddress(EthFlipAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.ValVoidSignature()),
		common.HexToHash("0x00000000000000000000000064d922894153be9eef7b7218dc565d1d0ce2a092"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e0c959c42b000000000000000000000000000000000000000000000000002386f26fc100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 16,
	TxHash:      common.HexToHash("0xc6ff19de9299e5b290ba2d52fdb4662360ca86376613d78ee546244866a0be2d"),
	TxIndex:     74,
	BlockHash:   fakes.FakeHash,
	Index:       75,
	Removed:     false,
}

var ValVoidEventLog = core.EventLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawValVoidLog,
	Transformed: false,
}

var valVoidModel = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.ValVoidTable,
	OrderedColumns: []event.ColumnName{
		event.HeaderFK,
		event.LogFK,
		event.AddressFK,
		constants.MsgSenderColumn,
	},
	ColumnValues: event.ColumnValues{
		event.HeaderFK: ValVoidEventLog.HeaderID,
		event.LogFK:    ValVoidEventLog.ID,
	},
}

func ValVoidModel() event.InsertionModel { return CopyModel(valVoidModel) }
