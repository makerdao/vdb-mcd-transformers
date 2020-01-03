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

var (
	frobData = "0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e0760887034554482d41000000000000000000000000000000000000000000000000000000000000000000000000000000eeec867b3f51ab5b619d582481bf53eea930b074000000000000000000000000eeec867b3f51ab5b619d582481bf53eea930b074000000000000000000000000eeec867b3f51ab5b619d582481bf53eea930b0740000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000016345785d8a000000000000000000000000000000000000000000000000000000000000"
)

var rawVatFrobLogWithPositiveDart = types.Log{
	Address: common.HexToAddress(VatAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.VatFrobSignature()),
		common.HexToHash("0x4554480000000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x000000000000000000000000eeec867b3f51ab5b619d582481bf53eea930b074"),
		common.HexToHash("0x000000000000000000000000fc7440e2ed4a3aeb14d40c00f02a14221be0474d"),
	},
	Data:        hexutil.MustDecode(frobData),
	BlockNumber: 10512592,
	TxHash:      common.HexToHash("0x10277b770bcd569cd3c943db2228153435ee1320eaab1f3a64fb8d5732d44c2e"),
	TxIndex:     123,
	BlockHash:   fakes.FakeHash,
	Index:       7,
	Removed:     false,
}

var VatFrobHeaderSyncLogWithPositiveDart = core.HeaderSyncLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawVatFrobLogWithPositiveDart,
	Transformed: false,
}

func VatFrobModelWithPositiveDart() event.InsertionModel {
	return CopyModel(vatFrobModelWithPositiveDart)
}

var vatFrobModelWithPositiveDart = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.VatFrobTable,
	OrderedColumns: []event.ColumnName{
		event.HeaderFK,
		event.LogFK,
		constants.UrnColumn,
		constants.VColumn,
		constants.WColumn,
		constants.DinkColumn,
		constants.DartColumn,
	},
	ColumnValues: event.ColumnValues{
		constants.VColumn:    "0xFc7440E2Ed4A3AEb14d40c00f02a14221Be0474d",
		constants.WColumn:    "0xEEec867B3F51ab5b619d582481BF53eea930b074",
		constants.DinkColumn: "0",
		constants.DartColumn: "100000000000000000",
		event.HeaderFK:       VatFrobHeaderSyncLogWithPositiveDart.HeaderID,
		event.LogFK:          VatFrobHeaderSyncLogWithPositiveDart.ID,
	},
}

var rawVatFrobLogWithNegativeDink = types.Log{
	Address: common.HexToAddress(VatAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.VatFrobSignature()),
		common.HexToHash("0x4554482d41000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x0000000000000000000000005c8c8e5895b9ccf34acf391c99e13c79ee2efb46"),
		common.HexToHash("0x000000000000000000000000fc7440e2ed4a3aeb14d40c00f02a14221be0474d"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e0760887034554482d410000000000000000000000000000000000000000000000000000000000000000000000000000005c8c8e5895b9ccf34acf391c99e13c79ee2efb46000000000000000000000000fc7440e2ed4a3aeb14d40c00f02a14221be0474d0000000000000000000000005c8c8e5895b9ccf34acf391c99e13c79ee2efb46ffffffffffffffffffffffffffffffffffffffffffffffffffe3940ad9cc0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 10713692,
	TxHash:      common.HexToHash("0x45bbf5d06e13d9b149f906226b4e64c98ef2132130633fa27a9e8f51fbedf4e2"),
	TxIndex:     0,
	BlockHash:   common.HexToHash("0x999bed34f207adb3e9860588bd6baa9f54414e55e43a933d88396271be45633d"),
	Index:       7,
	Removed:     false,
}

var VatFrobHeaderSyncLogWithNegativeDink = core.HeaderSyncLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawVatFrobLogWithNegativeDink,
	Transformed: false,
}

func VatFrobModelWithNegativeDink() event.InsertionModel {
	return CopyModel(vatFrobModelWithNegativeDink)
}

var vatFrobModelWithNegativeDink = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.VatFrobTable,
	OrderedColumns: []event.ColumnName{
		event.HeaderFK,
		event.LogFK,
		constants.UrnColumn,
		constants.VColumn,
		constants.WColumn,
		constants.DinkColumn,
		constants.DartColumn,
	},
	ColumnValues: event.ColumnValues{
		constants.VColumn:    "0xFc7440E2Ed4A3AEb14d40c00f02a14221Be0474d",
		constants.WColumn:    "0x5c8c8e5895B9cCf34ACF391C99E13C79EE2eFb46",
		constants.DinkColumn: "-8000000000000000",
		constants.DartColumn: "0",
		event.HeaderFK:       VatFrobHeaderSyncLogWithNegativeDink.HeaderID,
		event.LogFK:          VatFrobHeaderSyncLogWithNegativeDink.ID,
	},
}
