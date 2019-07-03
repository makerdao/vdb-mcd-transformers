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
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/transformers/events/flip_kick"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

var (
	flipID          = big.NewInt(1)
	lot             = big.NewInt(10)
	bid             = big.NewInt(25)
	tab             = big.NewInt(50)
	usr             = "0x7340e006f4135BA6970D43bf43d88DCAD4e7a8CA"
	gal             = "0x07Fa9eF6609cA7921112231F8f195138ebbA2977"
	contractAddress = constants.EthFlipContractAddressA()
	rawLog, _       = json.Marshal(EthFlipKickLog)
)

var (
	flipKickTransactionHash = "0xd11ab35cfb1ad71f790d3dd488cc1a2046080e765b150e8997aa0200947d4a9b"
	flipKickData            = "0x0000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000190000000000000000000000000000000000000000000000000000000000000032"
	FlipKickBlockNumber     = int64(10)
)

var EthFlipKickLog = types.Log{
	Address: common.HexToAddress(constants.EthFlipContractAddressA()),
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

var FlipKickEntity = flip_kick.FlipKickEntity{
	Id:               flipID,
	Lot:              lot,
	Bid:              bid,
	Tab:              tab,
	Usr:              common.HexToAddress(usr),
	Gal:              common.HexToAddress(gal),
	ContractAddress:  common.HexToAddress(contractAddress),
	TransactionIndex: EthFlipKickLog.TxIndex,
	LogIndex:         EthFlipKickLog.Index,
	Raw:              EthFlipKickLog,
}

var FlipKickModel = flip_kick.FlipKickModel{
	BidId:            flipID.String(),
	Lot:              lot.String(),
	Bid:              bid.String(),
	Tab:              tab.String(),
	Usr:              usr,
	Gal:              gal,
	ContractAddress:  contractAddress,
	TransactionIndex: EthFlipKickLog.TxIndex,
	LogIndex:         EthFlipKickLog.Index,
	Raw:              rawLog,
}

type FlipKickDBRow struct {
	ID       int64
	HeaderId int64 `db:"header_id"`
	flip_kick.FlipKickModel
}
