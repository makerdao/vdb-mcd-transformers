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

	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"

	"github.com/makerdao/vulcanizedb/pkg/core"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vulcanizedb/pkg/fakes"

	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
)

var rawVatMoveLog = types.Log{
	Address: common.HexToAddress(VatAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.VatMoveSignature()),
		common.HexToHash("0x000000000000000000000000a730d1ff8b6bc74a26d54c20a9dda539909bab0e"),
		common.HexToHash("0x000000000000000000000000b730d1ff8b6bc74a26d54c20a9dda539909bab0e"),
		common.HexToHash("0x000000000000000000000000000000000000000000000000000000000000002a"),
	},
	Data:        hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000006478f19470a730d1ff8b6bc74a26d54c20a9dda539909bab0e000000000000000000000000b730d1ff8b6bc74a26d54c20a9dda539909bab0e000000000000000000000000000000000000000000000000000000000000000000000000000000000000002a"),
	BlockNumber: 10,
	TxHash:      common.HexToHash("0xe8f39fbb7fea3621f543868f19b1114e305aff6a063a30d32835ff1012526f91"),
	TxIndex:     7,
	BlockHash:   fakes.FakeHash,
	Index:       8,
	Removed:     false,
}

var VatMoveHeaderSyncLog = core.HeaderSyncLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawVatMoveLog,
	Transformed: false,
}

var VatMoveModel = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.VatMoveTable,
	OrderedColumns: []event.ColumnName{
		event.HeaderFK, constants.SrcColumn, constants.DstColumn, constants.RadColumn, event.LogFK,
	},
	ColumnValues: event.ColumnValues{
		constants.SrcColumn: "0xA730d1FF8B6Bc74a26d54c20a9dda539909BaB0e",
		constants.DstColumn: "0xB730D1fF8b6BC74a26D54c20a9ddA539909BAb0e",
		constants.RadColumn: "42",
		event.HeaderFK:      VatMoveHeaderSyncLog.HeaderID,
		event.LogFK:         VatMoveHeaderSyncLog.ID,
	},
}
