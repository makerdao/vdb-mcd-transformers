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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
)

var rawCatFileChopLog = types.Log{
	Address: common.HexToAddress(CatAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.CatFileChopLumpSignature()),
		common.HexToHash("0x000000000000000000000000dc984d513a0f9ca9aa602d4df8517677918936e3"),
		common.HexToHash("0x434f4c342d410000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x63686f7000000000000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e01a0b287e434f4c342d41000000000000000000000000000000000000000000000000000063686f70000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000033b2e3c9fd0803ce800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 110,
	TxHash:      common.HexToHash("0xe32dfe6afd7ea28475569756fc30f0eea6ad4cfd32f67436ff1d1c805e4382df"),
	TxIndex:     13,
	BlockHash:   fakes.FakeHash,
	Index:       1,
	Removed:     false,
}

var CatFileChopEventLog = core.EventLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawCatFileChopLog,
	Transformed: false,
}

func CatFileChopModel() event.InsertionModel { return CopyModel(catFileChopModel) }

var catFileChopModel = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.CatFileChopLumpTable,
	OrderedColumns: []event.ColumnName{
		event.HeaderFK,
		event.LogFK,
		constants.MsgSenderColumn,
		constants.IlkColumn,
		constants.WhatColumn,
		constants.DataColumn,
	},
	ColumnValues: event.ColumnValues{
		event.HeaderFK: CatFileChopEventLog.HeaderID,
		event.LogFK:    CatFileChopEventLog.ID,
		// msg.sender
		// constants.IlkIdColumn changes with each insertion so it's best to assert directly in the test itself
		constants.WhatColumn: "chop",
		constants.DataColumn: "1000000000000000000000000000",
	},
}

var rawCatFileLumpLog = types.Log{
	Address: common.HexToAddress(CatAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.CatFileChopLumpSignature()),
		common.HexToHash("0x000000000000000000000000dc984d513a0f9ca9aa602d4df8517677918936e3"),
		common.HexToHash("0x434f4c342d410000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x6c756d7000000000000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e01a0b287e434f4c342d4100000000000000000000000000000000000000000000000000006c756d7000000000000000000000000000000000000000000000000000000000000000000000000000000006d79f82328ea3da61e066ebb2f88a00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 10980157,
	TxHash:      common.HexToHash("0xd8ff700a91f7216fd0d019e3f097ca581068fc0ef0dd4ace6eab6476df6a1987"),
	TxIndex:     15,
	BlockHash:   fakes.FakeHash,
	Index:       3,
	Removed:     false,
}

var CatFileLumpEventLog = core.EventLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawCatFileLumpLog,
	Transformed: false,
}

func CatFileLumpModel() event.InsertionModel { return CopyModel(catFileLumpModel) }

var catFileLumpModel = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.CatFileChopLumpTable,
	OrderedColumns: []event.ColumnName{
		event.HeaderFK,
		event.LogFK,
		constants.MsgSenderColumn,
		constants.IlkColumn,
		constants.WhatColumn,
		constants.DataColumn,
	},
	ColumnValues: event.ColumnValues{
		event.HeaderFK: CatFileLumpEventLog.HeaderID,
		event.LogFK:    CatFileLumpEventLog.ID,
		// msg.sender
		// constants.IlkIdColumn changes with each insertion so it's best to assert directly in the test itself
		constants.WhatColumn: "lump",
		constants.DataColumn: "10000000000000000000000000000000000000000000000000",
	},
}

var rawCatFileFlipLog = types.Log{
	Address: common.HexToAddress(CatAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.CatFileFlipSignature()),
		common.HexToHash("0x000000000000000000000000dc984d513a0f9ca9aa602d4df8517677918936e3"),
		common.HexToHash("0x434f4c312d410000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x666c697000000000000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e0ebecb39d434f4c312d410000000000000000000000000000000000000000000000000000666c6970000000000000000000000000000000000000000000000000000000000000000000000000000000006e8032435c84b08e30f27bfbb812ee365a095b3100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 10980092,
	TxHash:      common.HexToHash("0x11ee9c04fb4f68ea02a55cd4e67502d3f0ed19d45e0e5ec919f5981990d0f69e"),
	TxIndex:     0,
	BlockHash:   fakes.FakeHash,
	Index:       2,
	Removed:     false,
}

var CatFileFlipEventLog = core.EventLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawCatFileFlipLog,
	Transformed: false,
}

func CatFileFlipModel() event.InsertionModel { return CopyModel(catFileFlipModel) }

var catFileFlipModel = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.CatFileFlipTable,
	OrderedColumns: []event.ColumnName{
		event.HeaderFK,
		constants.IlkColumn,
		constants.WhatColumn,
		constants.FlipColumn,
		event.LogFK,
	},
	ColumnValues: event.ColumnValues{
		event.HeaderFK:       CatFileFlipEventLog.HeaderID,
		constants.WhatColumn: "flip",
		constants.FlipColumn: "0x6E8032435c84B08E30F27bfbb812Ee365A095b31",
		event.LogFK:          CatFileFlipEventLog.ID,
	},
}

var rawCatFileVowLog = types.Log{
	Address: common.HexToAddress(CatAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.CatFileVowSignature()),
		common.HexToHash("0x0000000000000000000000003652c2af10cbbdb753c3b46489db5226b73e6497"),
		common.HexToHash("0x766f770000000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x00000000000000000000000017560834075da3db54f737db74377e799c865821"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000044d4e8be83766f77000000000000000000000000000000000000000000000000000000000000000000000000000000000017560834075da3db54f737db74377e799c86582100000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 87,
	TxHash:      common.HexToHash("0x6515c7dfe53f0ad83ce1173fa99032c24a07cfd8b5d5a1c1f80486c99dd52800"),
	TxIndex:     11,
	BlockHash:   fakes.FakeHash,
	Index:       2,
	Removed:     false,
}

var CatFileVowEventLog = core.EventLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawCatFileVowLog,
	Transformed: false,
}

func CatFileVowModel() event.InsertionModel { return CopyModel(catFileVowModel) }

var catFileVowModel = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.CatFileVowTable,
	OrderedColumns: []event.ColumnName{
		event.HeaderFK,
		event.LogFK,
		constants.MsgSenderColumn,
		constants.WhatColumn,
		constants.DataColumn,
	},
	ColumnValues: event.ColumnValues{
		event.HeaderFK:       CatFileVowEventLog.HeaderID,
		constants.WhatColumn: "vow",
		constants.DataColumn: "0x17560834075DA3Db54f737db74377E799c865821",
		event.LogFK:          CatFileVowEventLog.ID,
	},
}
