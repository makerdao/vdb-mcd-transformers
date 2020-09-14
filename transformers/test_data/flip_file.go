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
		common.HexToHash("0x0000000000000000000000000000000000000000000000000e4b4b8af6a70000"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e029ae811462656700000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e4b4b8af6a70000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 10769102,
	TxHash:      common.HexToHash("0x9125c2c8795a0a872386e1ceda091c603db94f737a4557e784a96352086bd985"),
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
