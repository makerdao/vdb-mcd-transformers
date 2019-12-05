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
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"math/rand"

	"github.com/makerdao/vulcanizedb/pkg/core"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vulcanizedb/pkg/fakes"

	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
)

var rawVatInitLog = types.Log{
	Address: common.HexToAddress(VatAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.VatInitSignature()),
		common.HexToHash("0x66616b6520696c6b000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000243b66319566616b6520696c6b000000000000000000000000000000000000000000000000"),
	BlockNumber: 24,
	TxHash:      common.HexToHash("0xe8f39fbb7fea3621f543868f19b1114e305aff6a063a30d32835ff1012526f91"),
	TxIndex:     7,
	BlockHash:   fakes.FakeHash,
	Index:       8,
	Removed:     false,
}

var VatInitHeaderSyncLog = core.HeaderSyncLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawVatInitLog,
	Transformed: false,
}

var VatInitModel = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.VatInitTable,
	OrderedColumns: []event.ColumnName{
		event.HeaderFK, constants.IlkColumn, event.LogFK,
	},
	ColumnValues: event.ColumnValues{
		event.HeaderFK: VatInitHeaderSyncLog.HeaderID,
		event.LogFK:    VatInitHeaderSyncLog.ID,
		//constants.IlkColumn
	},
}
