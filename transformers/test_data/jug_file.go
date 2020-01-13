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

var rawJugFileIlkLog = types.Log{
	Address: common.HexToAddress(JugAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.JugFileIlkSignature()),
		common.HexToHash("0x000000000000000000000000127232c33f9b051e3703294de3c1e03e15f8a33f"),
		common.HexToHash("0x434f4c322d410000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x6475747900000000000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e01a0b287e434f4c322d41000000000000000000000000000000000000000000000000000064757479000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000033b2e3cacd278c7503e82c100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 10980334,
	TxHash:      common.HexToHash("0xa38f6bb83ae8c1bd239c883e3553e71d712db77bb3954851cc6ed5468a821613"),
	TxIndex:     2,
	BlockHash:   fakes.FakeHash,
	Index:       1,
	Removed:     false,
}

var JugFileIlkEventLog = core.EventLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawJugFileIlkLog,
	Transformed: false,
}

var jugFileIlkModel = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.JugFileIlkTable,
	OrderedColumns: []event.ColumnName{
		event.HeaderFK, constants.IlkColumn, constants.WhatColumn, constants.DataColumn, event.LogFK,
	},
	ColumnValues: event.ColumnValues{
		constants.WhatColumn: "duty",
		constants.DataColumn: "1000000000937303470807876289",
		// Ilk ID
		event.HeaderFK: JugFileIlkEventLog.HeaderID,
		event.LogFK:    JugFileIlkEventLog.ID,
	},
}

func JugFileIlkModel() event.InsertionModel { return CopyModel(jugFileIlkModel) }

var rawJugFileBaseLog = types.Log{
	Address: common.HexToAddress(JugAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.JugFileBaseSignature()),
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

var JugFileBaseEventLog = core.EventLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawJugFileBaseLog,
	Transformed: false,
}

var jugFileBaseModel = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.JugFileBaseTable,
	OrderedColumns: []event.ColumnName{
		event.HeaderFK, constants.WhatColumn, constants.DataColumn, event.LogFK,
	},
	ColumnValues: event.ColumnValues{
		constants.WhatColumn: "fake what",
		constants.DataColumn: big.NewInt(123).String(),
		event.HeaderFK:       JugFileBaseEventLog.HeaderID,
		event.LogFK:          JugFileBaseEventLog.ID,
	},
}

func JugFileBaseModel() event.InsertionModel { return CopyModel(jugFileBaseModel) }

var rawJugFileVowLog = types.Log{
	Address: common.HexToAddress(JugAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.JugFileVowSignature()),
		common.HexToHash("0x0000000000000000000000003652c2af10cbbdb753c3b46489db5226b73e6497"),
		common.HexToHash("0x766f770000000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x00000000000000000000000017560834075da3db54f737db74377e799c865821"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000044e9b674b9766f77000000000000000000000000000000000000000000000000000000000017560834075da3db54f737db74377e799c86582100000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 51,
	TxHash:      common.HexToHash("0x586e26b71b41fcd6905044dbe8f0cca300517542278f74a9b925c4f800fed85c"),
	TxIndex:     14,
	BlockHash:   fakes.FakeHash,
	Index:       17,
	Removed:     false,
}

var JugFileVowEventLog = core.EventLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawJugFileVowLog,
	Transformed: false,
}

var jugFileVowModel = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.JugFileVowTable,
	OrderedColumns: []event.ColumnName{
		event.HeaderFK, constants.WhatColumn, constants.DataColumn, event.LogFK,
	},
	ColumnValues: event.ColumnValues{
		constants.WhatColumn: "vow",
		constants.DataColumn: "0x17560834075DA3Db54f737db74377E799c865821",
		event.HeaderFK:       JugFileVowEventLog.HeaderID,
		event.LogFK:          JugFileVowEventLog.ID,
	},
}

func JugFileVowModel() event.InsertionModel { return CopyModel(jugFileVowModel) }
