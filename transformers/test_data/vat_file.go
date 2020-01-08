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

var rawVatFileDebtCeilingLog = types.Log{
	Address: common.HexToAddress(VatAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.VatFileDebtCeilingSignature()),
		common.HexToHash("0x4c696e6500000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x0000000000000000000002ac3a4edbbfb8014e3ba83411e915e8000000000000"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000004429ae81144c696e65000000000000000000000000000000000000000000000000000000000000000000000000000002ac3a4edbbfb8014e3ba83411e915e800000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 10577169,
	TxHash:      common.HexToHash("0x0ec18121d45f96293d9d759fd7564db4186a1aa69552f3106dd1afdeffdc9106"),
	TxIndex:     333,
	BlockHash:   fakes.FakeHash,
	Index:       15,
	Removed:     false,
}

var VatFileDebtCeilingEventLog = core.EventLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawVatFileDebtCeilingLog,
	Transformed: false,
}

func VatFileDebtCeilingModel() event.InsertionModel { return CopyModel(vatFileDebtCeilingModel) }

var vatFileDebtCeilingModel = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.VatFileDebtCeilingTable,
	OrderedColumns: []event.ColumnName{
		event.HeaderFK, constants.WhatColumn, constants.DataColumn, event.LogFK,
	},
	ColumnValues: event.ColumnValues{
		constants.WhatColumn: "Line",
		constants.DataColumn: "1000000000000000000000000000000000000000000000000000",
		event.HeaderFK:       VatFileDebtCeilingEventLog.HeaderID,
		event.LogFK:          VatFileDebtCeilingEventLog.ID,
	},
}

var rawVatFileIlkDustLog = types.Log{
	Address: common.HexToAddress(VatAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.VatFileIlkSignature()),                                       //sig
		common.HexToHash("0x5245500000000000000000000000000000000000000000000000000000000000"),  //ilk
		common.HexToHash("0x6475737400000000000000000000000000000000000000000000000000000000"),  //what
		common.HexToHash("0x000000000000000000000000000000000000000832600000bee4f14727b555555"), //val
	},

	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000641a0b287e52455000000000000000000000000000000000000000000000000000000000006475737400000000000000000000000000000000000000000000000000000000000000000000000000000000000000832600000bee4f14727b55555500000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 12,
	TxHash:      common.HexToHash("0x2e27c962a697d4f7ec5d3206d0c058bd510f7593a711f082e55f3b62d44d8dd9"),
	TxIndex:     112,
	BlockHash:   fakes.FakeHash,
	Index:       15,
	Removed:     false,
}

var VatFileIlkDustEventLog = core.EventLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawVatFileIlkDustLog,
	Transformed: false,
}

func VatFileIlkDustModel() event.InsertionModel { return CopyModel(vatFileIlkDustModel) }

var vatFileIlkDustModel = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.VatFileIlkTable,
	OrderedColumns: []event.ColumnName{
		event.HeaderFK, event.LogFK, constants.IlkColumn, constants.WhatColumn, constants.DataColumn,
	},
	ColumnValues: event.ColumnValues{
		constants.WhatColumn: "dust",
		constants.DataColumn: "10390649719961925488562719249749",
		event.HeaderFK:       VatFileIlkDustEventLog.HeaderID,
		event.LogFK:          VatFileIlkDustEventLog.ID,
	},
}

var rawVatFileIlkLineLog = types.Log{
	Address: common.HexToAddress(VatAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.VatFileIlkSignature()),
		common.HexToHash("0x5245500000000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x6c696e6500000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x00000000000000000000000000000000000000000000000bee4f14727b555555"),
	},

	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000641a0b287e52455000000000000000000000000000000000000000000000000000000000006c696e6500000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000bee4f14727b55555500000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 12,
	TxHash:      common.HexToHash("0x2e27c962a697d4f7ec5d3206d0c058bd510f7593a711f082e55f3b62d44d8dd9"),
	TxIndex:     112,
	BlockHash:   fakes.FakeHash,
	Index:       15,
	Removed:     false,
}

var VatFileIlkLineEventLog = core.EventLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawVatFileIlkLineLog,
	Transformed: false,
}

func VatFileIlkLineModel() event.InsertionModel { return CopyModel(vatFileIlkLineModel) }

var vatFileIlkLineModel = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.VatFileIlkTable,
	OrderedColumns: []event.ColumnName{
		event.HeaderFK, event.LogFK, constants.IlkColumn, constants.WhatColumn, constants.DataColumn,
	},
	ColumnValues: event.ColumnValues{
		constants.WhatColumn: "line",
		constants.DataColumn: "220086151196920075605",
		event.HeaderFK:       VatFileIlkLineEventLog.HeaderID,
		event.LogFK:          VatFileIlkLineEventLog.ID,
	},
}

var rawVatFileIlkSpotLog = types.Log{
	Address: common.HexToAddress(VatAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.VatFileIlkSignature()),
		common.HexToHash("0x4554480000000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x73706f7400000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x00000000000000000000000000000000000000012714e40bee4f14727b555555"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000641a0b287e455448000000000000000000000000000000000000000000000000000000000073706f740000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000012714e40bee4f14727b55555500000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 10501158,
	TxHash:      common.HexToHash("0x657baea11882f6a73b0088382fa9b9b7ba84f0c1183af909d93ab0fe2d10c292"),
	TxIndex:     111,
	BlockHash:   fakes.FakeHash,
	Index:       14,
	Removed:     false,
}

var VatFileIlkSpotEventLog = core.EventLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawVatFileIlkSpotLog,
	Transformed: false,
}

func VatFileIlkSpotModel() event.InsertionModel { return CopyModel(vatFileIlkSpotModel) }

var vatFileIlkSpotModel = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.VatFileIlkTable,
	OrderedColumns: []event.ColumnName{
		event.HeaderFK, event.LogFK, constants.IlkColumn, constants.WhatColumn, constants.DataColumn,
	},
	ColumnValues: event.ColumnValues{
		constants.WhatColumn: "spot",
		constants.DataColumn: "91323333333333333333333333333",
		event.HeaderFK:       VatFileIlkSpotEventLog.HeaderID,
		event.LogFK:          VatFileIlkSpotEventLog.ID,
	},
}
