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
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

var EthVowFlogLog = types.Log{
	Address: common.HexToAddress(constants.VowContractAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.VowFlogSignature()),
		common.HexToHash("0x0000000000000000000000008e84a1e068d77059cbe263c43ad0cdc130863313"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000539"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000002435aee16f0000000000000000000000000000000000000000000000000000000000000539"),
	BlockNumber: 11,
	TxHash:      common.HexToHash("0x47ffd75c1cda1d5c08219755743663a3790e4f5ae9e1fcb85bb3fe0d74bb7109"),
	TxIndex:     4,
	BlockHash:   fakes.FakeHash,
	Index:       3,
	Removed:     false,
}

var rawVowFlogLog, _ = json.Marshal(EthVowFlogLog)
var VowFlogModel = shared.InsertionModel{
	TableName: "vow_flog",
	OrderedColumns: []string{
		"header_id", "era", "log_idx", "tx_idx", "raw_log",
	},
	ColumnValues: shared.ColumnValues{
		"era":     "1337",
		"log_idx": EthVowFlogLog.Index,
		"tx_idx":  EthVowFlogLog.TxIndex,
		"raw_log": rawVowFlogLog,
	},
	ForeignKeyValues: shared.ForeignKeyValues{},
}
