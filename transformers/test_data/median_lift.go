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

var rawMedianLiftLogWithFourAccounts = types.Log{
	Address: common.HexToAddress(EthMedianAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.MedianLiftSignature()),
		common.HexToHash("0x000000000000000000000000c45e7858eef1318337a803ede8c5a9be12e2b40f"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000020"),
		common.HexToHash("0x000000000000000000000000000000000000000000000000000000000000000a"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e0943181060000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000a0000000000000000000000006bdbc0ccc17d72a33bf72a4657781a37dc2aa94e00000000000000000000000026c45f7b0e456e36fc85781488a3cd42a57ccbd200000000000000000000000020c576f989ee94e571f027b30314acf709267f7c000000000000000000000000fcb1fb52e114b364b3aab63d8a6f65fe8dcbef9d000000000000000000000000c2de180006ed15273f8dc59c436b954b"),
	BlockNumber: 8936530,
	TxHash:      common.HexToHash("0xd17875c308e4778ebe4335e445d84e9b280e181cd60e65ecce68f43da5e390b8"),
	TxIndex:     10,
	BlockHash:   fakes.FakeHash,
	Index:       11,
	Removed:     false,
}

var MedianLiftLogWithFourAccounts = core.EventLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawMedianLiftLogWithFourAccounts,
	Transformed: false,
}

func MedianLiftModelWithFourAccounts() event.InsertionModel {
	return CopyModel(medianLiftModelWithFourAccounts)
}

var medianLiftModelWithFourAccounts = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.MedianLiftTable,
	OrderedColumns: []event.ColumnName{
		event.HeaderFK,
		event.LogFK,
		event.AddressFK,
		constants.MsgSenderColumn,
		constants.A0Column,
		constants.A1Column,
		constants.A2Column,
		constants.A3Column,
	},
	ColumnValues: event.ColumnValues{
		event.HeaderFK: MedianLiftLogWithFourAccounts.HeaderID,
		event.LogFK:    MedianLiftLogWithFourAccounts.ID,
	},
}

var rawMedianLiftLogWithOneAccount = types.Log{
	Address: common.HexToAddress(EthMedianAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.MedianLiftSignature()),
		common.HexToHash("0x000000000000000000000000c45e7858eef1318337a803ede8c5a9be12e2b40f"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000020"),
		common.HexToHash("0x000000000000000000000000000000000000000000000000000000000000000a"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e08ef5eaf000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000001000000000000000000000000acb48fd097f1e0b24d3853bead826e5e9278b70000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 8936534,
	TxHash:      common.HexToHash("0x27f7834f778ec7d4289cf3337f8e428785c6d023164c02fc44565dbf2e26c49a"),
	TxIndex:     12,
	BlockHash:   fakes.FakeHash,
	Index:       14,
	Removed:     false,
}

var MedianLiftLogWithOneAccount = core.EventLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawMedianLiftLogWithOneAccount,
	Transformed: false,
}

func MedianLiftModelWithOneAccount() event.InsertionModel {
	return CopyModel(medianLiftModelWithOneAccount)
}

var medianLiftModelWithOneAccount = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.MedianLiftTable,
	OrderedColumns: []event.ColumnName{
		event.HeaderFK,
		event.LogFK,
		event.AddressFK,
		constants.MsgSenderColumn,
		constants.A0Column,
		constants.A1Column,
		constants.A2Column,
		constants.A3Column,
	},
	ColumnValues: event.ColumnValues{
		event.HeaderFK: MedianLiftLogWithOneAccount.HeaderID,
		event.LogFK:    MedianLiftLogWithOneAccount.ID,
	},
}
