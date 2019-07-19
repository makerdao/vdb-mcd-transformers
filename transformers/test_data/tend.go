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
	tendBidId           = int64(10)
	tendLot             = "8500000000000"
	tendBid             = "100000000000"
	tendData            = "0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e04b43ed12000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000007bb0f7b0800000000000000000000000000000000000000000000000000000000174876e80000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
	tendTransactionHash = "0xaa12e00846ceda4bf8ed33b1513c1972038c5152f8ca621dcb396652da9559b8"
)

var TendLogNote = types.Log{
	Address: common.HexToAddress(constants.GetContractAddress("MCD_FLAP")),
	Topics: []common.Hash{
		common.HexToHash(constants.TendSignature()),
		common.HexToHash("0x0000000000000000000000003a673843d27d037b206bb05eb1abbc7288d95e66"),
		common.HexToHash("0x000000000000000000000000000000000000000000000000000000000000000a"),
		common.HexToHash("0x000000000000000000000000000000000000000000000000000007bb0f7b0800"),
	},
	Data:        hexutil.MustDecode(tendData),
	BlockNumber: 11,
	TxHash:      common.HexToHash(tendTransactionHash),
	TxIndex:     10,
	BlockHash:   fakes.FakeHash,
	Index:       1,
	Removed:     false,
}

var rawTendLog, _ = json.Marshal(TendLogNote)
var TendModel = shared.InsertionModel{
	TableName: "tend",
	OrderedColumns: []string{
		"header_id", "bid_id", "lot", "bid", "contract_address", "log_idx", "tx_idx", "raw_log",
	},
	ColumnValues: shared.ColumnValues{
		"bid_id":           strconv.FormatInt(tendBidId, 10),
		"lot":              tendLot,
		"bid":              tendBid,
		"contract_address": constants.GetContractAddress("MCD_FLAP"),
		"log_idx":          TendLogNote.Index,
		"tx_idx":           TendLogNote.TxIndex,
		"raw_log":          rawTendLog,
	},
	ForeignKeyValues: shared.ForeignKeyValues{},
}
