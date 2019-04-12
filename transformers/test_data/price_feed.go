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
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/mcd_transformers/transformers/events/price_feeds"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

var (
	medianizerAddress = common.HexToAddress("0xefa5f53c62531cb29b8a8e298687a422b8793d72")
	blockNumber       = uint64(10653439)
	txIndex           = uint(2)
)

// https://kovan.etherscan.io/tx/0x7d144771f01d5b6b18ab777fb3b72ac7d61e48847cd7716c38d99d80c879942f
var EthPriceFeedLog = types.Log{
	Address:     medianizerAddress,
	Topics:      []common.Hash{common.HexToHash(KovanLogMedianPriceSignature)},
	Data:        common.FromHex("0x000000000000000000000000000000000000000000000000e6c7928c7e3f4c00000000000000000000000000000000000000000000000000000000005c9f06dc"),
	BlockNumber: blockNumber,
	TxHash:      common.HexToHash("0x7d144771f01d5b6b18ab777fb3b72ac7d61e48847cd7716c38d99d80c879942f"),
	TxIndex:     txIndex,
	BlockHash:   fakes.FakeHash,
	Index:       8,
	Removed:     false,
}

var rawPriceFeedLog, _ = json.Marshal(EthPriceFeedLog)
var PriceFeedModel = price_feeds.PriceFeedModel{
	BlockNumber:       blockNumber,
	MedianizerAddress: EthPriceFeedLog.Address.String(),
	UsdValue:          "16629421281200000000",
	Age:               "1553925852",
	LogIndex:          EthPriceFeedLog.Index,
	TransactionIndex:  EthPriceFeedLog.TxIndex,
	Raw:               rawPriceFeedLog,
}
