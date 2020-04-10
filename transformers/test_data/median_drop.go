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

var rawMedianDropLogWithFiveAccounts = types.Log{
	Address: common.HexToAddress(EthMedianAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.MedianDropSignature()),
		common.HexToHash("0x000000000000000000000000c45e7858eef1318337a803ede8c5a9be12e2b40f"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000020"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000005"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e046d4577d00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000005000000000000000000000000a52f23a651d1fa7c2610753c768103ee8c498f22000000000000000000000000ce91db32ad1c91278a56cbb2d8f24f9315043de90000000000000000000000003482f7a06db71f8ecac04f882546a66081311667000000000000000000000000702f365e1e559d9dc7b1af6ab9be64feb9c4d822000000000000000000000000ae37ab846ce92cf19031e602bf7dd3ae"),
	BlockNumber: 8936530,
	TxHash:      common.HexToHash("0x27f7834f778ec7d4289cf3337f8e428785c6d023164c02fc44565dbf2e26c49a"),
	TxIndex:     10,
	BlockHash:   fakes.FakeHash,
	Index:       11,
	Removed:     false,
}

var MedianDropLogWithFiveAccounts = core.EventLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawMedianDropLogWithFiveAccounts,
	Transformed: false,
}

func MedianDropModelWithFiveAccounts() event.InsertionModel {
	return CopyModel(medianDropModelWithFiveAccounts)
}

var medianDropModelWithFiveAccounts = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.MedianDropTable,
	OrderedColumns: []event.ColumnName{
		event.HeaderFK,
		event.LogFK,
		event.AddressFK,
		constants.MsgSenderColumn,
		constants.ALengthColumn,
		constants.AColumn,
	},
	ColumnValues: event.ColumnValues{
		event.HeaderFK:          MedianDropLogWithFiveAccounts.HeaderID,
		event.LogFK:             MedianDropLogWithFiveAccounts.ID,
		constants.ALengthColumn: "5",
	},
}

var rawMedianDropLogWithOneAccount = types.Log{
	Address: common.HexToAddress(EthMedianAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.MedianDropSignature()),
		common.HexToHash("0x000000000000000000000000c45e7858eef1318337a803ede8c5a9be12e2b40f"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000020"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e08ef5eaf000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000001000000000000000000000000acb48fd097f1e0b24d3853bead826e5e9278b70000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 8936534,
	TxHash:      common.HexToHash("0x27f7834f778ec7d4289cf3337f8e428785c6d023164c02fc44565dbf2e26c49a"),
	TxIndex:     12,
	BlockHash:   fakes.FakeHash,
	Index:       14,
	Removed:     false,
}

var MedianDropLogWithOneAccount = core.EventLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawMedianDropLogWithOneAccount,
	Transformed: false,
}

func MedianDropModelWithOneAccount() event.InsertionModel {
	return CopyModel(medianDropModelWithOneAccount)
}

var medianDropModelWithOneAccount = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.MedianDropTable,
	OrderedColumns: []event.ColumnName{
		event.HeaderFK,
		event.LogFK,
		event.AddressFK,
		constants.MsgSenderColumn,
		constants.ALengthColumn,
		constants.AColumn,
	},
	ColumnValues: event.ColumnValues{
		event.HeaderFK:          MedianDropLogWithOneAccount.HeaderID,
		event.LogFK:             MedianDropLogWithOneAccount.ID,
		constants.ALengthColumn: "1",
	},
}
