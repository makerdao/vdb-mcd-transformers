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
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/big"
	"math/rand"

	"github.com/vulcanize/mcd_transformers/transformers/events/flip_kick"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

var (
	flipID          = big.NewInt(1)
	lot             = big.NewInt(10)
	bid             = big.NewInt(25)
	tab             = big.NewInt(50)
	FakeUrn         = "0x7340e006f4135BA6970D43bf43d88DCAD4e7a8CA"
	gal             = "0x07Fa9eF6609cA7921112231F8f195138ebbA2977"
	contractAddress = EthFlipAddress()
	rawLog, _       = json.Marshal(FlipKickHeaderSyncLog)
)

var (
	flipKickTransactionHash = "0xd11ab35cfb1ad71f790d3dd488cc1a2046080e765b150e8997aa0200947d4a9b"
	flipKickData            = "0x0000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000190000000000000000000000000000000000000000000000000000000000000032"
	FlipKickBlockNumber     = int64(10)
)

var rawFlipKickLog = types.Log{
	Address: common.HexToAddress(EthFlipAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.FlipKickSignature()),
		common.HexToHash("0x0000000000000000000000007340e006f4135ba6970d43bf43d88dcad4e7a8ca"),
		common.HexToHash("0x00000000000000000000000007fa9ef6609ca7921112231f8f195138ebba2977"),
	},
	Data:        hexutil.MustDecode(flipKickData),
	BlockNumber: uint64(FlipKickBlockNumber),
	TxHash:      common.HexToHash(flipKickTransactionHash),
	TxIndex:     999,
	BlockHash:   fakes.FakeHash,
	Index:       1,
	Removed:     false,
}

var FlipKickHeaderSyncLog = core.HeaderSyncLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawFlipKickLog,
	Transformed: false,
}

var FlipKickEntity = flip_kick.FlipKickEntity{
	Id:              flipID,
	Lot:             lot,
	Bid:             bid,
	Tab:             tab,
	Usr:             common.HexToAddress(FakeUrn),
	Gal:             common.HexToAddress(gal),
	ContractAddress: common.HexToAddress(contractAddress),
	HeaderID:        FlipKickHeaderSyncLog.HeaderID,
	LogID:           FlipKickHeaderSyncLog.ID,
}

var FlipKickModel = shared.InsertionModel{
	SchemaName: "maker",
	TableName:  "flip_kick",
	OrderedColumns: []string{
		constants.HeaderFK, constants.LogFK, "bid_id", "lot", "bid", "tab", "usr", "gal", "address_id",
	},
	ColumnValues: shared.ColumnValues{
		constants.HeaderFK: FlipKickEntity.HeaderID,
		constants.LogFK:    FlipKickEntity.LogID,
		"bid_id":           flipID.String(),
		"lot":              lot.String(),
		"bid":              bid.String(),
		"tab":              tab.String(),
		"usr":              FakeUrn,
		"gal":              gal,
	},
	ForeignKeyValues: shared.ForeignKeyValues{
		constants.AddressFK: contractAddress,
	},
}

type FlipKickDBRow struct {
	ID int64
	flip_kick.FlipKickModel
}
