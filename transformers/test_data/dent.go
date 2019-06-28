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
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/transformers/events/dent"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

var (
	dentData            = "0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e05ff3a382000000000000000000000000000000000000000000000000002386f26fc1000000000000000000000000000000000000000000000000000000470de4df820000000000000000000000000000000000000000000000000000006a94d74f43000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
	dentTransactionHash = "0x5a210319fcd31eea5959fedb4a1b20881c21a21976e23ff19dff3b44cc1c71e8"
	dentBidId           = "10000000000000000"
	dentLot             = "20000000000000000"
	dentBid             = "30000000000000000"
	dentGuy             = "0x64d922894153BE9EEf7b7218dc565d1D0Ce2a092"
)

var EthDentLog = types.Log{
	Address: common.HexToAddress(constants.OldFlipperContractAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.DentSignature()),
		common.HexToHash("0x00000000000000000000000064d922894153be9eef7b7218dc565d1d0ce2a092"),
		common.HexToHash("0x000000000000000000000000000000000000000000000000002386f26fc10000"),
		common.HexToHash("0x00000000000000000000000000000000000000000000000000470de4df820000"),
	},
	Data:        hexutil.MustDecode(dentData),
	BlockNumber: 15,
	TxHash:      common.HexToHash(dentTransactionHash),
	TxIndex:     5,
	BlockHash:   fakes.FakeHash,
	Index:       2,
	Removed:     false,
}

var dentRawJson, _ = json.Marshal(EthDentLog)
var DentModel = dent.DentModel{
	BidId:            dentBidId,
	Lot:              dentLot,
	Bid:              dentBid,
	ContractAddress:  EthDentLog.Address.Hex(),
	LogIndex:         EthDentLog.Index,
	TransactionIndex: EthDentLog.TxIndex,
	Raw:              dentRawJson,
}
