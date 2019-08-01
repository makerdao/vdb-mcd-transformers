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

var EthVatInitLog = types.Log{
	Address: common.HexToAddress(VatAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.VatInitSignature()),
		common.HexToHash("0x66616b6520696c6b000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000243b66319566616b6520696c6b000000000000000000000000000000000000000000000000"),
	BlockNumber: 24,
	TxHash:      common.HexToHash("0xe8f39fbb7fea3621f543868f19b1114e305aff6a063a30d32835ff1012526f91"),
	TxIndex:     7,
	BlockHash:   fakes.FakeHash,
	Index:       8,
	Removed:     false,
}

var rawVatInitLog, _ = json.Marshal(EthVatInitLog)
var VatInitModel = shared.InsertionModel{
	TableName: "vat_init",
	OrderedColumns: []string{
		"header_id", string(constants.IlkFK), "log_idx", "tx_idx", "raw_log",
	},
	ColumnValues: shared.ColumnValues{
		"log_idx": EthVatInitLog.Index,
		"tx_idx":  EthVatInitLog.TxIndex,
		"raw_log": rawVatInitLog,
	},
	ForeignKeyValues: shared.ForeignKeyValues{
		constants.IlkFK: "0x66616b6520696c6b000000000000000000000000000000000000000000000000",
	},
}
