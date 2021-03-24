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
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
)

var rawEthSpotPokeLog = types.Log{
	Address: common.HexToAddress(SpotAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.SpotPokeSignature()),
	},
	Data: hexutil.MustDecode("0x434f4c352d410000000000000000000000000000000000000000000000000000" +
		"000000000000000000000000000000000000000000000000013c6d703eec3700" +
		"00000000000000000000000000000000000000000026c69aba83d25409ffca1a",
	),
	BlockNumber: 11257491,
	TxHash:      common.HexToHash("0x1103dd21f6c43d6f895d41935918119d35d000a109f9353c4575959ba01206bd"),
	TxIndex:     1,
	BlockHash:   common.HexToHash("0x94c99a3b2c77f7cbc401f010c50cd298546cd1084c4d436dde4d22a8b7bbe7e1"),
	Index:       2,
	Removed:     false,
}
var SpotPokeEventLog = core.EventLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawEthSpotPokeLog,
	Transformed: false,
}

func SpotPokeModel() event.InsertionModel { return CopyModel(spotPokeModel) }

const SpotPokeIlkHex = "0x434f4c352d410000000000000000000000000000000000000000000000000000"

var spotPokeModel = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.SpotPokeTable,
	OrderedColumns: []event.ColumnName{
		event.HeaderFK,
		event.LogFK,
		constants.IlkColumn,
		constants.ValueColumn,
		constants.SpotColumn,
	},
	ColumnValues: event.ColumnValues{
		event.HeaderFK:        SpotPokeEventLog.HeaderID,
		event.LogFK:           SpotPokeEventLog.ID,
		constants.ValueColumn: "89066421500000000.000000",
		constants.SpotColumn:  "46877063947368421052631578",
	},
}
