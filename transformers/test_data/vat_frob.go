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
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_frob"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

var (
	frobData = "0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e0760887034554482d41000000000000000000000000000000000000000000000000000000000000000000000000000000eeec867b3f51ab5b619d582481bf53eea930b074000000000000000000000000eeec867b3f51ab5b619d582481bf53eea930b074000000000000000000000000eeec867b3f51ab5b619d582481bf53eea930b0740000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000016345785d8a000000000000000000000000000000000000000000000000000000000000"
)

var EthVatFrobLog = types.Log{
	Address: common.HexToAddress(KovanVatContractAddress),
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

var rawVatFrobLog, _ = json.Marshal(EthVatFrobLog)
var VatFrobModel = vat_frob.VatFrobModel{
	Ilk:              "4554480000000000000000000000000000000000000000000000000000000000",
	Urn:              "0xEEec867B3F51ab5b619d582481BF53eea930b074",
	V:                "0xFc7440E2Ed4A3AEb14d40c00f02a14221Be0474d",
	W:                "0xEEec867B3F51ab5b619d582481BF53eea930b074",
	Dink:             "0",
	Dart:             "100000000000000000",
	LogIndex:         EthVatFrobLog.Index,
	TransactionIndex: EthVatFrobLog.TxIndex,
	Raw:              rawVatFrobLog,
}
