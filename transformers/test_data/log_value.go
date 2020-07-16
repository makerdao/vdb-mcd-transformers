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
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/log_value"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
)

const (
	TemporaryLogValueBlockNumber = int64(14824113)
	logValueData                 = "0x000000000000000000000000000000000000000000000009f19a45256de70000"
	TemporaryLogValueTransaction = "0x81c446bceeaeb9b5117cda43eeb5926e63d4057247018bd0c6cc54c07c8eb15b"
)

var logValueVal, _ = new(big.Int).SetString("183430000000000000000", 10)

var rawLogValueLog = types.Log{
	Address:     common.HexToAddress(OsmEthAddress()),
	Topics:      []common.Hash{common.HexToHash(constants.LogValueSignature())},
	Data:        hexutil.MustDecode(logValueData),
	BlockNumber: uint64(TemporaryLogValueBlockNumber),
	TxHash:      common.HexToHash(TemporaryLogValueTransaction),
	TxIndex:     0,
	BlockHash:   fakes.FakeHash,
	Index:       7,
	Removed:     false,
}
var LogValueEventLog = core.EventLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawLogValueLog,
	Transformed: false,
}

func LogValueModel() event.InsertionModel { return CopyModel(logValueModel) }

var logValueModel = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.LogValueTable,
	OrderedColumns: []event.ColumnName{
		event.HeaderFK, event.LogFK, log_value.Val,
	},
	ColumnValues: event.ColumnValues{
		event.HeaderFK: LogValueEventLog.HeaderID,
		event.LogFK:    LogValueEventLog.ID,
		log_value.Val:  logValueVal.String(),
	},
}
