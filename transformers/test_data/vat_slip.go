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

	"github.com/vulcanize/mcd_transformers/transformers/events/vat_slip"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

var EthVatSlipLog = types.Log{
	Address: common.HexToAddress(KovanVatContractAddress),
	Topics: []common.Hash{
		common.HexToHash(KovanVatSlipSignature),
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

var rawVatSlipLog, _ = json.Marshal(EthVatSlipLog)
var VatSlipModel = vat_slip.VatSlipModel{
	Ilk:              "4554482d41000000000000000000000000000000000000000000000000000000",
	Usr:              "0xFc7440E2Ed4A3AEb14d40c00f02a14221Be0474d",
	Wad:              "115792089237316195423570985008687907853269984665640564039457579007913129639936",
	TransactionIndex: EthVatSlipLog.TxIndex,
	LogIndex:         EthVatSlipLog.Index,
	Raw:              rawVatSlipLog,
}
