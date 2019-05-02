// VulcanizeDB
// Copyright © 2018 Vulcanize

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

	"github.com/vulcanize/mcd_transformers/transformers/events/vat_frob"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

var (
	frobData = "0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e0760887034554482d41000000000000000000000000000000000000000000000000000000000000000000000000000000eeec867b3f51ab5b619d582481bf53eea930b074000000000000000000000000eeec867b3f51ab5b619d582481bf53eea930b074000000000000000000000000eeec867b3f51ab5b619d582481bf53eea930b0740000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000016345785d8a000000000000000000000000000000000000000000000000000000000000"
)

var EthVatFrobLogWithPositiveDart = types.Log{
	Address: common.HexToAddress(constants.VatContractAddress()),
	Topics: []common.Hash{
		common.HexToHash(KovanVatFrobSignature),
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

var rawVatFrobLogWithPositiveDart, _ = json.Marshal(EthVatFrobLogWithPositiveDart)
var VatFrobModelWithPositiveDart = vat_frob.VatFrobModel{
	Ilk:              "4554480000000000000000000000000000000000000000000000000000000000",
	Urn:              "0xEEec867B3F51ab5b619d582481BF53eea930b074",
	V:                "0xFc7440E2Ed4A3AEb14d40c00f02a14221Be0474d",
	W:                "0xEEec867B3F51ab5b619d582481BF53eea930b074",
	Dink:             "0",
	Dart:             "100000000000000000",
	LogIndex:         EthVatFrobLogWithPositiveDart.Index,
	TransactionIndex: EthVatFrobLogWithPositiveDart.TxIndex,
	Raw:              rawVatFrobLogWithPositiveDart,
}

var EthVatFrobLogWithNegativeDink = types.Log{
	Address: common.HexToAddress(constants.VatContractAddress()),
	Topics: []common.Hash{
		common.HexToHash(KovanVatFrobSignature),
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

var rawVatFrobLogWithNegativeDink, _ = json.Marshal(EthVatFrobLogWithNegativeDink)
var VatFrobModelWithNegativeDink = vat_frob.VatFrobModel{
	Ilk:              "4554482d41000000000000000000000000000000000000000000000000000000",
	Urn:              "0x5c8c8e5895B9cCf34ACF391C99E13C79EE2eFb46",
	V:                "0xFc7440E2Ed4A3AEb14d40c00f02a14221Be0474d",
	W:                "0x5c8c8e5895B9cCf34ACF391C99E13C79EE2eFb46",
	Dink:             "-8000000000000000",
	Dart:             "0",
	LogIndex:         7,
	TransactionIndex: 0,
	Raw:              rawVatFrobLogWithNegativeDink,
}
