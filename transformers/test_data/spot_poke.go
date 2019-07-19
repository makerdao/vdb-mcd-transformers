//  VulcanizeDB
//  Copyright © 2019 Vulcanize
//
//  This program is free software: you can redistribute it and/or modify
//  it under the terms of the GNU Affero General Public License as published by
//  the Free Software Foundation, either version 3 of the License, or
//  (at your option) any later version.
//
//  This program is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//  GNU Affero General Public License for more details.
//
//  You should have received a copy of the GNU Affero General Public License
//  along with this program.  If not, see <http://www.gnu.org/licenses/>.

package test_data

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/mcd_transformers/transformers/events/spot_poke"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"math/big"
)

var EthSpotPokeLog = types.Log{
	Address: common.HexToAddress(constants.GetContractAddress("MCD_SPOT")),
	Topics: []common.Hash{
		common.HexToHash(constants.SpotPokeSignature()),
	},
	Data:        hexutil.MustDecode("0x434f4c352d410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000013c6d703eec370000000000000000000000000000000000000000000026c69aba83d25409ffca1a"),
	BlockNumber: 11257491,
	TxHash:      common.HexToHash("0x1103dd21f6c43d6f895d41935918119d35d000a109f9353c4575959ba01206bd"),
	TxIndex:     1,
	BlockHash:   common.HexToHash("0x94c99a3b2c77f7cbc401f010c50cd298546cd1084c4d436dde4d22a8b7bbe7e1"),
	Index:       2,
	Removed:     false,
}

func spot() *big.Int {
	spot := big.Int{}
	spot.SetString("46877063947368421052631578", 10)
	return &spot
}

var SpotPokeEntity = spot_poke.SpotPokeEntity{
	Ilk:              [32]byte{67, 79, 76, 53, 45, 65, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	Val:              [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 60, 109, 112, 62, 236, 55, 0},
	Spot:             spot(),
	TransactionIndex: 1,
	LogIndex:         2,
	Raw:              EthSpotPokeLog,
}

var rawLogJson, _ = json.Marshal(EthSpotPokeLog)
var SpotPokeModel = spot_poke.SpotPokeModel{
	Ilk:              "0x434f4c352d410000000000000000000000000000000000000000000000000000",
	Value:            "89066421500000000.000000",
	Spot:             "46877063947368421052631578",
	TransactionIndex: 1,
	LogIndex:         2,
	Raw:              rawLogJson,
}
