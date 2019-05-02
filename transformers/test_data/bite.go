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

	"github.com/vulcanize/mcd_transformers/transformers/events/bite"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

const (
	TemporaryBiteBlockNumber = int64(26)
	TemporaryBiteData        = "0x00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000300000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000005"
	TemporaryBiteTransaction = "0x5c698f13940a2153440c6d19660878bc90219d9298fdcf37365aa8d88d40fc42"
)

var (
	biteInk        = big.NewInt(1)
	biteArt        = big.NewInt(2)
	biteTab        = big.NewInt(3)
	biteFlip       = big.NewInt(4)
	biteRawJson, _ = json.Marshal(EthBiteLog)
	biteIlk        = [32]byte{69, 84, 72, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	biteUrn        = common.BytesToAddress([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 216, 180, 20, 126, 218, 128, 254, 199, 18, 42, 225, 109, 162, 71, 156, 189, 127, 251})
)

var EthBiteLog = types.Log{
	Address: common.HexToAddress(constants.CatContractAddress()),
	Topics: []common.Hash{
		common.HexToHash(KovanBiteSignature),
		common.HexToHash("0x4554480000000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x0000000000000000000000000000d8b4147eda80fec7122ae16da2479cbd7ffb"),
	},
	Data:        hexutil.MustDecode(TemporaryBiteData),
	BlockNumber: uint64(TemporaryBiteBlockNumber),
	TxHash:      common.HexToHash(TemporaryBiteTransaction),
	TxIndex:     111,
	BlockHash:   fakes.FakeHash,
	Index:       7,
	Removed:     false,
}

var BiteEntity = bite.BiteEntity{
	Ilk:              biteIlk,
	Urn:              biteUrn,
	Ink:              biteInk,
	Art:              biteArt,
	Tab:              biteTab,
	Flip:             biteFlip,
	LogIndex:         EthBiteLog.Index,
	TransactionIndex: EthBiteLog.TxIndex,
	Raw:              EthBiteLog,
}

var BiteModel = bite.BiteModel{
	Ilk:              "4554480000000000000000000000000000000000000000000000000000000000",
	Urn:              "0000d8b4147eda80fec7122ae16da2479cbd7ffb",
	Ink:              biteInk.String(),
	Art:              biteArt.String(),
	Tab:              biteTab.String(),
	Flip:             biteFlip.String(),
	LogIndex:         EthBiteLog.Index,
	TransactionIndex: EthBiteLog.TxIndex,
	Raw:              biteRawJson,
}
