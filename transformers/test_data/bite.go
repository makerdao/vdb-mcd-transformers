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
	biteData                 = "0x00000000000000000000000000000000000000000000000000000002540be40000000000000000000000000000000000000000000000000000000004a817c80000000000000000000000000000000000000000000000000000000006fc23ac000000000000000000000000007d7bee5fcfd8028cf7b00876c5b1421c800561a600000000000000000000000000000000000000000000000000000009502f9000"
	TemporaryBiteTransaction = "0x5c698f13940a2153440c6d19660878bc90219d9298fdcf37365aa8d88d40fc42"
)

var (
	biteInk        = big.NewInt(10000000000)
	biteArt        = big.NewInt(20000000000)
	biteTab        = big.NewInt(30000000000)
	biteFlip       = common.HexToAddress("0x7d7bEe5fCfD8028cf7b00876C5b1421c800561A6")
	biteID         = big.NewInt(40000000000)
	biteRawJson, _ = json.Marshal(EthBiteLog)
	biteIlk        = [32]byte{69, 84, 72, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	biteUrn        = common.HexToAddress("0x7d7bEe5fCfD8028cf7b00876C5b1421c800561A6")
)

var EthBiteLog = types.Log{
	Address: common.HexToAddress(CatAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.BiteSignature()),
		common.HexToHash("0x4554480000000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x0000000000000000000000007d7bee5fcfd8028cf7b00876c5b1421c800561a6"),
	},
	Data:        hexutil.MustDecode(biteData),
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
	Id:               biteID,
	LogIndex:         EthBiteLog.Index,
	TransactionIndex: EthBiteLog.TxIndex,
	Raw:              EthBiteLog,
}

var BiteModel = bite.BiteModel{
	Ilk:              "0x4554480000000000000000000000000000000000000000000000000000000000",
	Urn:              "0x7d7bEe5fCfD8028cf7b00876C5b1421c800561A6",
	Ink:              biteInk.String(),
	Art:              biteArt.String(),
	Tab:              biteTab.String(),
	Flip:             "0x7d7bEe5fCfD8028cf7b00876C5b1421c800561A6",
	Id:               biteID.String(),
	LogIndex:         EthBiteLog.Index,
	TransactionIndex: EthBiteLog.TxIndex,
	Raw:              biteRawJson,
}
