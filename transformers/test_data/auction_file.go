package test_data

import (
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
)

var rawAuctionFileLog = types.Log{
	Address: common.HexToAddress(FlipUsdcBV107Address()),
	Topics: []common.Hash{
		common.HexToHash(constants.AuctionFileSignature()),
		common.HexToHash("0x000000000000000000000000be8e3e3618f7474f8cb1d074a26affef007e98fb"),
		common.HexToHash("0x6265670000000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000e4b4b8af6a70000"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e029ae811462656700000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e4b4b8af6a70000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 10201136,
	TxHash:      common.HexToHash("0x5f7242e46336b03bf798f67414e7b6eeea4518135a2cafbfafc86e8bdd631110"),
	TxIndex:     76,
	BlockHash:   common.HexToHash("0x81ef9ef9974a50195a977affa693e52aadc8cda0aa3a18f54912039c9e26bccf"),
	Index:       72,
	Removed:     false,
}

var AuctionFileEventLog = core.EventLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawAuctionFileLog,
	Transformed: false,
}

var auctionFileModel = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.AuctionFileTable,
	OrderedColumns: []event.ColumnName{
		event.HeaderFK,
		event.LogFK,
		event.AddressFK,
		constants.MsgSenderColumn,
		constants.WhatColumn,
		constants.DataColumn,
	},
	ColumnValues: event.ColumnValues{
		event.HeaderFK: AuctionFileEventLog.HeaderID,
		event.LogFK:    AuctionFileEventLog.ID,
		// Contract address fk
		// Msg sender fk
		constants.WhatColumn: "beg",
		constants.DataColumn: "1030000000000000000",
	},
}

func AuctionFileModel() event.InsertionModel { return CopyModel(auctionFileModel) }
