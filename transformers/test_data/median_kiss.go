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

var rawMedianKissSingleLog = types.Log{
	Address: common.HexToAddress(MedianEthAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.MedianKissSingleSignature()),
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

var MedianKissSingleLog = core.EventLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawMedianKissSingleLog,
	Transformed: false,
}

func MedianKissSingleModel() event.InsertionModel { return CopyModel(medianKissSingleModel) }

var medianKissSingleModel = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.MedianKissSingleTable,
	OrderedColumns: []event.ColumnName{
		event.HeaderFK,
		event.LogFK,
		event.AddressFK,
		constants.MsgSenderColumn,
		constants.AColumn,
	},
	ColumnValues: event.ColumnValues{
		event.HeaderFK: MedianKissSingleLog.HeaderID,
		event.LogFK:    MedianKissSingleLog.ID,
	},
}

var rawMedianKissBatchLogOneAddress = types.Log{
	Address: common.HexToAddress(MedianEthAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.MedianKissBatchSignature()),
		common.HexToHash("0x000000000000000000000000e87f55af91068a1da44095138f3d37c45894eb21"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000020"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e046d4577d00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000001000000000000000000000000a52f23a651d1fa7c2610753c768103ee8c498f2200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 8936530,
	TxHash:      common.HexToHash("0x27f7834f778ec7d4289cf3337f8e428785c6d023164c02fc44565dbf2e26c49a"),
	TxIndex:     10,
	BlockHash:   fakes.FakeHash,
	Index:       11,
	Removed:     false,
}

var MedianKissBatchLogOneAddress = core.EventLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawMedianKissBatchLogOneAddress,
	Transformed: false,
}

func MedianKissBatchModelOneAddress() event.InsertionModel {
	return CopyModel(medianKissBatchModelOneAddress)
}

var medianKissBatchModelOneAddress = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.MedianKissBatchTable,
	OrderedColumns: []event.ColumnName{
		event.HeaderFK,
		event.LogFK,
		event.AddressFK,
		constants.MsgSenderColumn,
		constants.ALengthColumn,
		constants.AColumn,
	},
	ColumnValues: event.ColumnValues{
		event.HeaderFK:          MedianKissBatchLogOneAddress.HeaderID,
		event.LogFK:             MedianKissBatchLogOneAddress.ID,
		constants.ALengthColumn: "1",
	},
}

var RawMedianKissBatchLogFiveAddresses = types.Log{
	Address: common.HexToAddress(MedianEthAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.MedianKissBatchSignature()),
		common.HexToHash("0x000000000000000000000000e87f55af91068a1da44095138f3d37c45894eb21"),
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

var MedianKissBatchLogFiveAddresses = core.EventLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         RawMedianKissBatchLogFiveAddresses,
	Transformed: false,
}

func MedianKissBatchModelFiveAddresses() event.InsertionModel {
	return CopyModel(medianKissBatchModelFiveAddresses)
}

var medianKissBatchModelFiveAddresses = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.MedianKissBatchTable,
	OrderedColumns: []event.ColumnName{
		event.HeaderFK,
		event.LogFK,
		event.AddressFK,
		constants.MsgSenderColumn,
		constants.ALengthColumn,
		constants.AColumn,
	},
	ColumnValues: event.ColumnValues{
		event.HeaderFK:          MedianKissBatchLogFiveAddresses.HeaderID,
		event.LogFK:             MedianKissBatchLogFiveAddresses.ID,
		constants.ALengthColumn: "5",
	},
}
