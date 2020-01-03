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

	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
)

var rawVatSuckLog = types.Log{
	Address: common.HexToAddress(VatAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.VatSuckSignature()),
		common.HexToHash("0x00000000000000000000000007fa9ef6609ca7921112231f8f195138ebba2977"),
		common.HexToHash("0x0000000000000000000000007340e006f4135ba6970d43bf43d88dcad4e7a8ca"),
		common.HexToHash("0x00000000000000000000000000000000000000000000003635c9adc5dea00000"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e0ee8cd74800000000000000000000000007fa9ef6609ca7921112231f8f195138ebba29770000000000000000000000007340e006f4135ba6970d43bf43d88dcad4e7a8ca00000000000000000000000000000000000000000000003635c9adc5dea0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 105,
	TxHash:      common.HexToHash("0x2730b707ef875c02ea45fd68f6d193320b85cf062b1860a02d1f1d407c845b65"),
	TxIndex:     2,
	BlockHash:   common.HexToHash("0x39185e33e15a6bd521240566bc3c5e34853ecd1af3212b000d50e7ca80d5cdbc"),
	Index:       3,
	Removed:     false,
}

var VatSuckHeaderSyncLog = core.HeaderSyncLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawVatSuckLog,
	Transformed: false,
}

var VatSuckModel = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.VatSuckTable,
	OrderedColumns: []event.ColumnName{
		event.HeaderFK, event.LogFK, constants.UColumn, constants.VColumn, constants.RadColumn,
	},
	ColumnValues: event.ColumnValues{
		constants.UColumn:   "0x07Fa9eF6609cA7921112231F8f195138ebbA2977",
		constants.VColumn:   "0x7340e006f4135BA6970D43bf43d88DCAD4e7a8CA",
		constants.RadColumn: "1000000000000000000000",
		event.HeaderFK:      VatSuckHeaderSyncLog.HeaderID,
		event.LogFK:         VatSuckHeaderSyncLog.ID,
	},
}
