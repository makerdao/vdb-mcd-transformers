// VulcanizeDB
// Copyright © 2019 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

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

var rawSpotFileMatLog = types.Log{
	Address: common.HexToAddress(SpotAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.SpotFileMatSignature()),
		common.HexToHash("0x00000000000000000000000071ce79fcfec71760d51f6b3589c0d9ec0ccb64a8"),
		common.HexToHash("0x4554482d41000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x6d61740000000000000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e01a0b287e4554482d410000000000000000000000000000000000000000000000000000006d61740000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004d8c55aefb8c05b5c00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 11257385,
	TxHash:      common.HexToHash("0xb4d19aaf5be5077db65aeeb16906a5498cfa94836952191258cc45966e1d7198"),
	TxIndex:     2,
	BlockHash:   common.HexToHash("0x968cd16acb356de42e9f3ab17583988b49173c0339af5afa3f79cecdbc111d69"),
	Index:       3,
	Removed:     false,
}

var SpotFileMatEventLog = core.EventLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawSpotFileMatLog,
	Transformed: false,
}

func SpotFileMatModel() event.InsertionModel { return CopyModel(spotFileMatModel) }

var spotFileMatModel = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.SpotFileMatTable,
	OrderedColumns: []event.ColumnName{
		event.HeaderFK,
		event.LogFK,
		constants.IlkColumn,
		constants.WhatColumn,
		constants.DataColumn,
		constants.MsgSenderColumn,
	},
	ColumnValues: event.ColumnValues{
		event.HeaderFK:       SpotFileMatEventLog.HeaderID,
		event.LogFK:          SpotFileMatEventLog.ID,
		constants.WhatColumn: "mat",
		constants.DataColumn: "1500000000000000000000000000",
		// msgSenderId
	},
}

var rawSpotFileParLog = types.Log{
	Address: common.HexToAddress(SpotAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.SpotFileParSignature()),
		common.HexToHash("0x00000000000000000000000064d922894153be9eef7b7218dc565d1d0ce2a092"),
		common.HexToHash("0x66616b6520776861740000000000000000000000000000000000000000000000"),
		common.HexToHash("0x000000000000000000000000000000000000000000000000000000000000007b"),
	},
	Data:        hexutil.MustDecode("0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000004429ae811466616b6520776861740000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000007b"),
	BlockNumber: 36,
	TxHash:      common.HexToHash("0xeeaa16de1d91c239b66773e8c2116a26cfeaaf5d962b31466c9bf047a5caa20f"),
	TxIndex:     13,
	BlockHash:   fakes.FakeHash,
	Index:       16,
	Removed:     false,
}

var SpotFileParEventLog = core.EventLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawSpotFileParLog,
	Transformed: false,
}

var spotFileParModel = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.SpotFileParTable,
	OrderedColumns: []event.ColumnName{
		event.HeaderFK, constants.WhatColumn, constants.DataColumn, event.LogFK, constants.MsgSenderColumn,
	},
	ColumnValues: event.ColumnValues{
		constants.WhatColumn: "fake what",
		constants.DataColumn: big.NewInt(123).String(),
		event.HeaderFK:       SpotFileParEventLog.HeaderID,
		event.LogFK:          SpotFileParEventLog.ID,
		// MsgSenderColumn
	},
}

func SpotFileParModel() event.InsertionModel { return CopyModel(spotFileParModel) }

var rawSpotFilePipLog = types.Log{
	Address: common.HexToAddress(SpotAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.SpotFilePipSignature()),
		common.HexToHash("0x00000000000000000000000073acbfb5b9413b0020164ee63dce4e1f71aba67c"),
		common.HexToHash("0x4554482d41000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x7069700000000000000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e0ebecb39d4554482d41000000000000000000000000000000000000000000000000000000706970000000000000000000000000000000000000000000000000000000000000000000000000000000000075dd74e8afe8110c8320ed397cccff3b8134d98100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 11257191,
	TxHash:      common.HexToHash("0xaae9e8bce346f86a01c5a3af137bc1f9bc7c0c767804a2b9b6356849aee0d7dd"),
	TxIndex:     1,
	BlockHash:   common.HexToHash("0xfa28e186578238fdd6b971add2ebe62a26dddf5ff971d50ee476c86b45362da1"),
	Index:       2,
	Removed:     false,
}

var SpotFilePipEventLog = core.EventLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawSpotFilePipLog,
	Transformed: false,
}

func SpotFilePipModel() event.InsertionModel { return CopyModel(spotFilePipModel) }

var spotFilePipModel = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.SpotFilePipTable,
	OrderedColumns: []event.ColumnName{
		event.HeaderFK,
		event.LogFK,
		constants.IlkColumn,
		constants.WhatColumn,
		constants.PipColumn,
	},
	ColumnValues: event.ColumnValues{
		event.HeaderFK:       SpotFilePipEventLog.HeaderID,
		event.LogFK:          SpotFilePipEventLog.ID,
		constants.WhatColumn: "pip",
		constants.PipColumn:  "0x75dD74e8afE8110C8320eD397CcCff3B8134d981",
	},
}
