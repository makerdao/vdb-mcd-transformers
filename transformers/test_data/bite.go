// VulcanizeDB
// Copyright © 2018 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//d GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package test_data

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/mcd_transformers/transformers/events/bite"
	"github.com/vulcanize/vulcanizedb/libraries/shared/factories/event"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/big"
	"math/rand"

	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

const (
	TemporaryBiteBlockNumber = int64(26)
	biteData                 = "0x00000000000000000000000000000000000000000000000000000002540be40000000000000000000000000000000000000000000000000000000004a817c80000000000000000000000000000000000000000000000000000000006fc23ac000000000000000000000000007d7bee5fcfd8028cf7b00876c5b1421c800561a600000000000000000000000000000000000000000000000000000009502f9000"
	TemporaryBiteTransaction = "0x5c698f13940a2153440c6d19660878bc90219d9298fdcf37365aa8d88d40fc42"
)

var (
	biteInk = big.NewInt(10000000000)
	biteArt = big.NewInt(20000000000)
	biteTab = big.NewInt(30000000000)
	biteID  = big.NewInt(40000000000)
)

var rawBiteLog = types.Log{
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
var BiteHeaderSyncLog = core.HeaderSyncLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawBiteLog,
	Transformed: false,
}

func BiteModel() event.InsertionModel { return CopyEventModel(biteModel) }

var biteModel = event.InsertionModel{
	SchemaName: "maker",
	TableName:  "bite",
	OrderedColumns: []event.ColumnName{
		constants.HeaderFK, constants.LogFK, constants.UrnColumn, bite.Ink, bite.Art, bite.Tab, bite.Flip, bite.Id,
	},
	ColumnValues: event.ColumnValues{
		constants.HeaderFK: BiteHeaderSyncLog.HeaderID,
		constants.LogFK:    BiteHeaderSyncLog.ID,
		// constants.UrnColumn: Can't assert against this since we don't know the ID...
		bite.Ink:  biteInk.String(),
		bite.Art:  biteArt.String(),
		bite.Tab:  biteTab.String(),
		bite.Flip: "0x7d7bEe5fCfD8028cf7b00876C5b1421c800561A6",
		bite.Id:   biteID.String(),
	},
}
