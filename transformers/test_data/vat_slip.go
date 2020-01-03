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

var rawVatSlipLogWithPositiveWad = types.Log{
	Address: common.HexToAddress(VatAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.VatSlipSignature()),
		common.HexToHash("0x4554482d41000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x0000000000000000000000005c8c8e5895b9ccf34acf391c99e13c79ee2efb46"),
		common.HexToHash("0x000000000000000000000000000000000000000000000000002386f26fc10000"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e07cdd3fde4554482d410000000000000000000000000000000000000000000000000000000000000000000000000000005c8c8e5895b9ccf34acf391c99e13c79ee2efb46000000000000000000000000000000000000000000000000002386f26fc1000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 10713689,
	TxHash:      common.HexToHash("0xf8a206ce1acb5c39125cab95456041afb4ccfbe496cf8850e982259128f5aafd"),
	TxIndex:     0,
	BlockHash:   common.HexToHash("0xde1338d81bd1c8e0472fa96e13d5fa58f7a215d499d8f17f15adbef7ef9586b8"),
	Index:       5,
	Removed:     false,
}

var VatSlipHeaderSyncLogWithPositiveWad = core.HeaderSyncLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawVatSlipLogWithPositiveWad,
	Transformed: false,
}

var vatSlipModelWithPositiveWad = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.VatSlipTable,
	OrderedColumns: []event.ColumnName{
		event.HeaderFK, event.LogFK, constants.IlkColumn, constants.UsrColumn, constants.WadColumn,
	},
	ColumnValues: event.ColumnValues{
		constants.UsrColumn: "0x5c8c8e5895B9cCf34ACF391C99E13C79EE2eFb46",
		constants.WadColumn: "10000000000000000",
		event.HeaderFK:      VatSlipHeaderSyncLogWithPositiveWad.HeaderID,
		event.LogFK:         VatSlipHeaderSyncLogWithPositiveWad.ID,
		constants.IlkColumn: "0x4554482d41000000000000000000000000000000000000000000000000000000",
	},
}

var rawVatSlipLogWithNegativeWad = types.Log{
	Address: common.HexToAddress(VatAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.VatSlipSignature()),
		common.HexToHash("0x4554482d41000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x000000000000000000000000fc7440e2ed4a3aeb14d40c00f02a14221be0474d"),
		common.HexToHash("0xffffffffffffffffffffffffffffffffffffffffffffffffffee3c86c81f8000"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e07cdd3fde4554482d41000000000000000000000000000000000000000000000000000000000000000000000000000000fc7440e2ed4a3aeb14d40c00f02a14221be0474dffffffffffffffffffffffffffffffffffffffffffffffffffee3c86c81f800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 10,
	TxHash:      common.HexToHash("0x2cb2c40a8385de94b05e47080216b2b10b7cfd45951aa06297f4e1d184e47118"),
	TxIndex:     3,
	BlockHash:   fakes.FakeHash,
	Index:       2,
	Removed:     false,
}

var VatSlipHeaderSyncLogWithNegativeWad = core.HeaderSyncLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawVatSlipLogWithNegativeWad,
	Transformed: false,
}

var vatSlipModelWithNegativeWad = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.VatSlipTable,
	OrderedColumns: []event.ColumnName{
		event.HeaderFK, event.LogFK, constants.IlkColumn, constants.UsrColumn, constants.WadColumn,
	},
	ColumnValues: event.ColumnValues{
		constants.UsrColumn: "0xFc7440E2Ed4A3AEb14d40c00f02a14221Be0474d",
		constants.WadColumn: "-5000000000000000",
		event.HeaderFK:      VatSlipHeaderSyncLogWithNegativeWad.HeaderID,
		event.LogFK:         VatSlipHeaderSyncLogWithNegativeWad.ID,
		//constants.IlkColumn DB state
	},
}

func VatSlipModelWithPositiveWad() event.InsertionModel {
	return CopyModel(vatSlipModelWithPositiveWad)
}
func VatSlipModelWithNegativeWad() event.InsertionModel {
	return CopyModel(vatSlipModelWithNegativeWad)
}
