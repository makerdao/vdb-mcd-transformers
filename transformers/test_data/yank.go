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

var EthYankLog = types.Log{
	Address: common.HexToAddress(constants.GetContractAddress("ETH_FLIP_A")),
	Topics: []common.Hash{
		common.HexToHash(constants.YankSignature()),
		common.HexToHash("0x00000000000000000000000064d922894153be9eef7b7218dc565d1d0ce2a092"),
		common.HexToHash("0x000000000000000000000000000000000000000000000000002386f26fc10000"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e026e027f1000000000000000000000000000000000000000000000000002386f26fc100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 106,
	TxHash:      common.HexToHash("0xd0735a21325a74f7dd312bf6b0e6e69ab58111f9023ac764810eafac587c19f0"),
	TxIndex:     3,
	BlockHash:   common.HexToHash("0x24adc827a634697a48bc19611f0008eaa4d19b6d58b262170aecc07c32e5da1f"),
	Index:       2,
	Removed:     false,
}

var rawYank, _ = json.Marshal(EthYankLog)
var YankModel = shared.InsertionModel{
	TableName: "yank",
	OrderedColumns: []string{
		"header_id", "bid_id", "contract_address", "log_idx", "tx_idx", "raw_log",
	},
	ColumnValues: shared.ColumnValues{
		"bid_id":           "10000000000000000",
		"contract_address": EthYankLog.Address.Hex(),
		"log_idx":          EthYankLog.Index,
		"tx_idx":           EthYankLog.TxIndex,
		"raw_log":          rawYank,
	},
	ForeignKeyValues: shared.ForeignKeyValues{},
}
