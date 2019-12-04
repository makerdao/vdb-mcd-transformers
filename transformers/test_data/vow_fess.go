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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
)

var rawVowFessLog = types.Log{
	Address: common.HexToAddress(VowAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.VowFessSignature()),
		common.HexToHash("0x0000000000000000000000002f34f22a00ee4b7a8f8bbc4eaee1658774c624e0"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000539"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x0000000000000000000000000000000000000000000002544faa778090e00000"),
	BlockNumber: 9377319,
	TxHash:      common.HexToHash("0x71f86b6154333d88811d151a9afecd00b39d6a326ef308ac97f66ca61264d7cb"),
	TxIndex:     4,
	BlockHash:   fakes.FakeHash,
	Index:       3,
	Removed:     false,
}

var VowFessHeaderSyncLog = core.HeaderSyncLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawVowFessLog,
	Transformed: false,
}

var VowFessModel = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.VowFessLabel,
	OrderedColumns: []event.ColumnName{
		event.HeaderFK, constants.TabColumn, event.LogFK,
	},
	ColumnValues: event.ColumnValues{
		constants.TabColumn: "1337",
		event.HeaderFK:      VowFessHeaderSyncLog.HeaderID,
		event.LogFK:         VowFessHeaderSyncLog.ID,
	},
}
