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
	"github.com/vulcanize/mcd_transformers/transformers/events/flop_kick"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/big"
)

var (
	EthFlopKickLog = types.Log{
		Address: common.HexToAddress(constants.GetContractAddress("MCD_FLOP")),
		Topics: []common.Hash{
			common.HexToHash(constants.FlopKickSignature()),
			common.HexToHash("0x0000000000000000000000007d7bee5fcfd8028cf7b00876c5b1421c800561a6"),
		},
		Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000006a94d74f4300000000000000000000000000000000000000000000000000000de0b6b3a76400000000000000000000000000000000000000000000000000001bc16d674ec80000"),
		BlockNumber: 19,
		TxHash:      common.HexToHash("0xd8fd67b37a6aa64a3cef4937204765183b180d8dc92eecd0d233f445526d31b5"),
		TxIndex:     uint(33),
		BlockHash:   fakes.FakeHash,
		Index:       32,
		Removed:     false,
	}

	FlopKickEntity = flop_kick.Entity{
		Id:               big.NewInt(30000000000000000),
		Lot:              big.NewInt(1000000000000000000),
		Bid:              big.NewInt(2000000000000000000),
		Gal:              common.HexToAddress("0x7d7bEe5fCfD8028cf7b00876C5b1421c800561A6"),
		TransactionIndex: EthFlopKickLog.TxIndex,
		LogIndex:         EthFlopKickLog.Index,
		Raw:              EthFlopKickLog,
	}

	rawFlopLogJson, _ = json.Marshal(EthFlopKickLog)
	FlopKickModel     = flop_kick.Model{
		BidId:            FlopKickEntity.Id.String(),
		Lot:              FlopKickEntity.Lot.String(),
		Bid:              FlopKickEntity.Bid.String(),
		Gal:              FlopKickEntity.Gal.Hex(),
		ContractAddress:  EthFlopKickLog.Address.Hex(),
		TransactionIndex: EthFlopKickLog.TxIndex,
		LogIndex:         EthFlopKickLog.Index,
		Raw:              rawFlopLogJson,
	}
)

type FlopKickDBResult struct {
	Id       int64
	HeaderId int64 `db:"header_id"`
	flop_kick.Model
}
