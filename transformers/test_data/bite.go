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
	"math/big"
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/fakes"

	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
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
	Address: common.HexToAddress(Cat100Address()),
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
var BiteEventLog = core.EventLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawBiteLog,
	Transformed: false,
}

func BiteModel() event.InsertionModel { return CopyModel(biteModel) }

var biteModel = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.BiteTable,
	OrderedColumns: []event.ColumnName{
		event.HeaderFK,
		event.LogFK,
		constants.UrnColumn,
		constants.InkColumn,
		constants.ArtColumn,
		constants.TabColumn,
		constants.FlipColumn,
		constants.BidIDColumn,
	},
	ColumnValues: event.ColumnValues{
		event.HeaderFK: BiteEventLog.HeaderID,
		event.LogFK:    BiteEventLog.ID,
		// constants.UrnColumn: Can't assert against this since we don't know the ID...
		constants.InkColumn:   biteInk.String(),
		constants.ArtColumn:   biteArt.String(),
		constants.TabColumn:   biteTab.String(),
		constants.FlipColumn:  "0x7d7bEe5fCfD8028cf7b00876C5b1421c800561A6",
		constants.BidIDColumn: biteID.String(),
	},
}
