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

var rawFlipFileCatLog = types.Log{
	Address: common.HexToAddress(FlipEthV110Address()),
	Topics: []common.Hash{
		common.HexToHash(constants.FlipFileCatSignature()),
		common.HexToHash("0x000000000000000000000000be8e3e3618f7474f8cb1d074a26affef007e98fb"),
		common.HexToHash("0x6361740000000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x000000000000000000000000a950524441892a31ebddf91d3ceefa04bf454466"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e029ae811462656700000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e4b4b8af6a70000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 8928165,
	TxHash:      common.HexToHash("0x00d3bfaa8425a58ec798663f36ffc3546dabe91b3e7994bc723585db50f3822d"),
	TxIndex:     999,
	BlockHash:   fakes.FakeHash,
	Index:       1,
	Removed:     false,
}

var FlipFileCatEventLog = core.EventLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawFlipFileCatLog,
	Transformed: false,
}

func FlipFileCatModel() event.InsertionModel { return CopyModel(flipFileCatModel) }

var flipFileCatModel = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.FlipFileCatTable,
	OrderedColumns: []event.ColumnName{
		event.HeaderFK,
		event.LogFK,
		event.AddressFK,
		constants.MsgSenderColumn,
		constants.WhatColumn,
		constants.DataColumn,
	},
	ColumnValues: event.ColumnValues{
		event.HeaderFK: FlipFileCatEventLog.HeaderID,
		event.LogFK:    FlipFileCatEventLog.ID,
		//AddressFK
		//MsgSender
	},
}
