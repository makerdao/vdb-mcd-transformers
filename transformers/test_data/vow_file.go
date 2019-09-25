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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"math/rand"
)

var rawVowFileLog = types.Log{
	Address: common.HexToAddress(VowAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.VowFileSignature()),
		common.HexToHash("0x00000000000000000000000064d922894153be9eef7b7218dc565d1d0ce2a092"),
		common.HexToHash("0x7761697400000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x00000000000000000000000000000000000000000000152d02c7e14af6800000"),
	},
	Data:        hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000004429ae8114776169740000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000152d02c7e14af680000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 10,
	TxHash:      common.Hash{},
	TxIndex:     11,
	BlockHash:   common.Hash{},
	Index:       12,
	Removed:     false,
}

var VowFileHeaderSyncLog = core.HeaderSyncLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawVowFileLog,
	Transformed: false,
}

var VowFileModel = shared.InsertionModel{
	SchemaName: "maker",
	TableName:  "vow_file",
	OrderedColumns: []string{
		constants.HeaderFK, "what", "data", constants.LogFK,
	},
	ColumnValues: shared.ColumnValues{
		"what":             "wait",
		"data":             "100000000000000000000000",
		constants.HeaderFK: VowFileHeaderSyncLog.HeaderID,
		constants.LogFK:    VowFileHeaderSyncLog.ID,
	},
	ForeignKeyValues: shared.ForeignKeyValues{},
}
