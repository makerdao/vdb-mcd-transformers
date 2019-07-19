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
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

var (
	flipTickBidId           = int64(10)
	flipTickData            = "0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e0fc7b6aee000000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
	flipTickTransactionHash = "0x6dc191fc774d5c5dc82bb292e6e2c4c62b5476b7fc9e589a89c3120448161966"
)

var FlipTickLogNote = types.Log{
	Address: common.HexToAddress(constants.GetContractAddress("ETH_FLIP_A")),
	Topics: []common.Hash{
		common.HexToHash(constants.FlipTickSignature()),
		common.HexToHash("0x000000000000000000000000da198bfdd2671d7ad4614c9cf2beb87cdfec1460"),
		common.HexToHash("0x000000000000000000000000000000000000000000000000000000000000000a"),
	},
	Data:        hexutil.MustDecode(flipTickData),
	BlockNumber: 11,
	TxHash:      common.HexToHash(flipTickTransactionHash),
	TxIndex:     10,
	BlockHash:   fakes.FakeHash,
	Index:       1,
	Removed:     false,
}

var rawFlipTickLog, _ = json.Marshal(FlipTickLogNote)
var FlipTickModel = shared.InsertionModel{
	TableName: "flip_tick",
	OrderedColumns: []string{
		"header_id", "bid_id", "contract_address", "log_idx", "tx_idx", "raw_log",
	},
	ColumnValues: shared.ColumnValues{
		"bid_id":           strconv.FormatInt(flipTickBidId, 10),
		"contract_address": FlipTickLogNote.Address.Hex(),
		"log_idx":          FlipTickLogNote.Index,
		"tx_idx":           FlipTickLogNote.TxIndex,
		"raw_log":          rawFlipTickLog,
	},
	ForeignKeyValues: shared.ForeignKeyValues{},
}
