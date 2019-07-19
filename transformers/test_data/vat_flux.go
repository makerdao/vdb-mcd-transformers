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

	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

var EthVatFluxLog = types.Log{
	Address: common.HexToAddress(constants.GetContractAddress("MCD_VAT")),
	Topics: []common.Hash{
		common.HexToHash(constants.VatFluxSignature()),
		common.HexToHash("0x66616b6520696c6b000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x00000000000000000000000007fa9ef6609ca7921112231f8f195138ebba2977"),
		common.HexToHash("0x0000000000000000000000007340e006f4135ba6970d43bf43d88dcad4e7a8ca"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e06111be2e66616b6520696c6b00000000000000000000000000000000000000000000000000000000000000000000000007fa9ef6609ca7921112231f8f195138ebba29770000000000000000000000007340e006f4135ba6970d43bf43d88dcad4e7a8ca000000000000000000000000000000000000000000000000000000e8d4a510000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 23,
	TxHash:      common.HexToHash("0x67db3532a08cb3ad3340eee79f4397d005cfbc9d721c1462018453f3af5f9286"),
	TxIndex:     115,
	BlockHash:   common.HexToHash("0x361c34aa03d509c47dc99deadd0678f0cf57d16d0143bba36d0bb7025bead343"),
	Index:       3,
	Removed:     false,
}

var rawFluxLog, _ = json.Marshal(EthVatFluxLog)
var VatFluxModel = shared.InsertionModel{
	TableName: "vat_flux",
	OrderedColumns: []string{
		"header_id", string(constants.IlkFK), "src", "dst", "wad", "tx_idx", "log_idx", "raw_log",
	},
	ColumnValues: shared.ColumnValues{
		"src":     "0x07Fa9eF6609cA7921112231F8f195138ebbA2977",
		"dst":     "0x7340e006f4135BA6970D43bf43d88DCAD4e7a8CA",
		"wad":     "1000000000000",
		"tx_idx":  EthVatFluxLog.TxIndex,
		"log_idx": EthVatFluxLog.Index,
		"raw_log": rawFluxLog,
	},
	ForeignKeyValues: shared.ForeignKeyValues{
		constants.IlkFK: "0x66616b6520696c6b000000000000000000000000000000000000000000000000",
	},
}
