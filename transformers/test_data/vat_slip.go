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
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

var EthVatSlipLogWithPositiveWad = types.Log{
	Address: common.HexToAddress(constants.GetContractAddress("MCD_VAT")),
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

var rawVatSlipLogWithPositiveWad, _ = json.Marshal(EthVatSlipLogWithPositiveWad)
var VatSlipModelWithPositiveWad = shared.InsertionModel{
	TableName: "vat_slip",
	OrderedColumns: []string{
		"header_id", string(constants.IlkFK), "usr", "wad", "tx_idx", "log_idx", "raw_log",
	},
	ColumnValues: shared.ColumnValues{
		"usr":     "0x5c8c8e5895B9cCf34ACF391C99E13C79EE2eFb46",
		"wad":     "10000000000000000",
		"tx_idx":  EthVatSlipLogWithPositiveWad.TxIndex,
		"log_idx": EthVatSlipLogWithPositiveWad.Index,
		"raw_log": rawVatSlipLogWithPositiveWad,
	},
	ForeignKeyValues: shared.ForeignKeyValues{
		constants.IlkFK: "0x4554482d41000000000000000000000000000000000000000000000000000000",
	},
}

var EthVatSlipLogWithNegativeWad = types.Log{
	Address: common.HexToAddress(constants.GetContractAddress("MCD_VAT")),
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

var rawVatSlipLogWithNegativeWad, _ = json.Marshal(EthVatSlipLogWithNegativeWad)
var VatSlipModelWithNegativeWad = shared.InsertionModel{
	TableName: "vat_slip",
	OrderedColumns: []string{
		"header_id", string(constants.IlkFK), "usr", "wad", "tx_idx", "log_idx", "raw_log",
	},
	ColumnValues: shared.ColumnValues{
		"usr":     "0xFc7440E2Ed4A3AEb14d40c00f02a14221Be0474d",
		"wad":     "-5000000000000000",
		"tx_idx":  EthVatSlipLogWithNegativeWad.TxIndex,
		"log_idx": EthVatSlipLogWithNegativeWad.Index,
		"raw_log": rawVatSlipLogWithNegativeWad,
	},
	ForeignKeyValues: shared.ForeignKeyValues{
		constants.IlkFK: "0x4554482d41000000000000000000000000000000000000000000000000000000",
	},
}
