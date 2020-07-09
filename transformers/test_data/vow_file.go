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
	"math/rand"

	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/pkg/core"
)

var rawVowFileLog = types.Log{
	Address: common.HexToAddress(VowAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.VowFileSignature()),
		common.HexToHash("0x00000000000000000000000064d922894153be9eef7b7218dc565d1d0ce2a092"),
		common.HexToHash("0x7761697400000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x00000000000000000000000000000000000000000000152d02c7e14af6800000"),
	},
	Data:        hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000004429ae8114776169740000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000152d02c7e14af680000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 10,
	TxHash:      common.Hash{},
	TxIndex:     11,
	BlockHash:   common.Hash{},
	Index:       12,
	Removed:     false,
}

var VowFileEventLog = core.EventLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawVowFileLog,
	Transformed: false,
}

var VowFileModel = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.VowFileTable,
	OrderedColumns: []event.ColumnName{
		event.HeaderFK, event.LogFK, constants.WhatColumn, constants.DataColumn,
	},
	ColumnValues: event.ColumnValues{
		constants.WhatColumn: "wait",
		constants.DataColumn: "100000000000000000000000",
		event.HeaderFK:       VowFileEventLog.HeaderID,
		event.LogFK:          VowFileEventLog.ID,
	},
}

var rawVowFileAuctionAddressLog = types.Log{
	Address: common.HexToAddress(VowAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.VowFileSignature()),
		common.HexToHash("0x000000000000000000000000be8e3e3618f7474f8cb1d074a26affef007e98fb"),
		common.HexToHash("0x666c6f7070657200000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x0000000000000000000000004d95a049d5b0b7d32058cd3f2163015747522e99"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e0d4e8be83666c6f70706572000000000000000000000000000000000000000000000000000000000000000000000000004d95a049d5b0b7d32058cd3f2163015747522e99000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 9017707,
	TxHash:      common.HexToHash("0x6fd3b8c2873015db136ae6ade1d7941642b9d2e8e6a70d2f624d48981d154be3"),
	TxIndex:     11,
	BlockHash:   common.Hash{},
	Index:       12,
	Removed:     false,
}

var VowFileAuctionAddressEventLog = core.EventLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawVowFileAuctionAddressLog,
	Transformed: false,
}

func VowFileAuctionAddressModel() event.InsertionModel { return vowFileAuctionAddressModel }

var vowFileAuctionAddressModel = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.VowFileAuctionAddressTable,
	OrderedColumns: []event.ColumnName{
		event.HeaderFK, event.LogFK, constants.WhatColumn, constants.DataColumn,
	},
	ColumnValues: event.ColumnValues{
		constants.WhatColumn: "flopper",
		constants.DataColumn: "0x4d95a049d5b0b7d32058cd3f2163015747522e99",
		event.HeaderFK:       VowFileAuctionAddressEventLog.HeaderID,
		event.LogFK:          VowFileAuctionAddressEventLog.ID,
	},
}
