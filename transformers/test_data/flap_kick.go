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
    "math/big"

    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/common/hexutil"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/vulcanize/vulcanizedb/pkg/fakes"

    "github.com/vulcanize/mcd_transformers/transformers/events/flap_kick"
    "github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

var EthFlapKickLog = types.Log{
	Address: common.HexToAddress(FlapperAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.FlapKickSignature()),
		common.HexToHash("0x0000000000000000000000007d7bee5fcfd8028cf7b00876c5b1421c800561a6"),
	},
	Data:        hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000246139ca80000000000000000000000000000000000000000000000000000de0b6b3a76400000000000000000000000000000000000000000000000000001bc16d674ec80000"),
	BlockNumber: 65,
	TxHash:      common.HexToHash("0xee7930b76b6e93974bd3f37824644ae42a89a3887a1131a7bcb3267ab4dc0169"),
	TxIndex:     66,
	BlockHash:   fakes.FakeHash,
	Index:       67,
	Removed:     false,
}

var FlapKickEntity = flap_kick.FlapKickEntity{
	Id:               big.NewInt(40000000000000),
	Lot:              big.NewInt(1000000000000000000),
	Bid:              big.NewInt(2000000000000000000),
	Gal:              common.HexToAddress("0x7d7bEe5fCfD8028cf7b00876C5b1421c800561A6"),
	Raw:              EthFlapKickLog,
	TransactionIndex: EthFlapKickLog.TxIndex,
	LogIndex:         EthFlapKickLog.Index,
}

var rawFlapKickLog, _ = json.Marshal(EthFlapKickLog)
var FlapKickModel = flap_kick.FlapKickModel{
	BidId:            FlapKickEntity.Id.String(),
	Lot:              FlapKickEntity.Lot.String(),
	Bid:              FlapKickEntity.Bid.String(),
	Gal:              FlapKickEntity.Gal.String(),
	ContractAddress:  EthFlapKickLog.Address.Hex(),
	LogIndex:         EthFlapKickLog.Index,
	TransactionIndex: EthFlapKickLog.TxIndex,
	Raw:              rawFlapKickLog,
}
