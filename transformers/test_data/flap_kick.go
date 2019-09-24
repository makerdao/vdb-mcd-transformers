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
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/transformers/events/flap_kick"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

var EthFlapKickLog = types.Log{
	Address:     common.HexToAddress(FlapAddress()),
	Topics:      []common.Hash{common.HexToHash(constants.FlapKickSignature())},
	Data:        hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000003b9aca000000000000000000000000000000000000000000000000000000000001312d00"),
	BlockNumber: 65,
	TxHash:      common.HexToHash("0xee7930b76b6e93974bd3f37824644ae42a89a3887a1131a7bcb3267ab4dc0169"),
	TxIndex:     66,
	BlockHash:   fakes.FakeHash,
	Index:       67,
	Removed:     false,
}

var FlapKickEntity = flap_kick.FlapKickEntity{
	Id:               big.NewInt(1),
	Lot:              big.NewInt(1000000000),
	Bid:              big.NewInt(20000000),
	Raw:              EthFlapKickLog,
	TransactionIndex: EthFlapKickLog.TxIndex,
	LogIndex:         EthFlapKickLog.Index,
}

var rawFlapKickLog, _ = json.Marshal(EthFlapKickLog)
var FlapKickModel = shared.InsertionModel{
	SchemaName: "maker",
	TableName:  "flap_kick",
	OrderedColumns: []string{
		"header_id", "bid_id", "lot", "bid", "address_id", "tx_idx", "log_idx", "raw_log",
	},
	ColumnValues: shared.ColumnValues{
		"bid_id":  FlapKickEntity.Id.String(),
		"lot":     FlapKickEntity.Lot.String(),
		"bid":     FlapKickEntity.Bid.String(),
		"log_idx": EthFlapKickLog.Index,
		"tx_idx":  EthFlapKickLog.TxIndex,
		"raw_log": rawFlapKickLog,
	},
	ForeignKeyValues: shared.ForeignKeyValues{
		constants.AddressFK: EthFlapKickLog.Address.Hex(),
	},
}
