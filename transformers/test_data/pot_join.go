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
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
)

var rawPotJoinLog = types.Log{
	Address: common.HexToAddress(PotAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.PotJoinSignature()),
		common.HexToHash("0x000000000000000000000000e7bc397dbd069fc7d0109c0636d06888bb50668c"),
		common.HexToHash("0x000000000000000000000000000000000000000000000000417fa3222791bd1a"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e0049878f3000000000000000000000000000000000000000000000000417fa3222791bd1a0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 105,
	TxHash:      common.HexToHash("0xe5bebbe1ea46d8b6d1515ce9ac4659f9e6774669c1f2009dcc7289c18b91e393"),
	TxIndex:     2,
	BlockHash:   common.HexToHash("0x3e011e52723db56476dc8cd45e1325f7bf3f3b2d89651253d6e8b66489f37d7c"),
	Index:       3,
	Removed:     false,
}

var PotJoinHeaderSyncLog = core.HeaderSyncLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawPotJoinLog,
	Transformed: false,
}

func PotJoinModel() event.InsertionModel { return CopyModel(potJoinModel) }

var potJoinModel = event.InsertionModel{
	SchemaName:     constants.MakerSchema,
	TableName:      constants.PotJoinTable,
	OrderedColumns: []event.ColumnName{event.HeaderFK, event.LogFK, constants.WadColumn},
	ColumnValues: event.ColumnValues{
		constants.WadColumn: "4719670301595647258",
		event.HeaderFK:      PotJoinHeaderSyncLog.HeaderID,
		event.LogFK:         PotJoinHeaderSyncLog.ID,
	},
}
