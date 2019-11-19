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
	"github.com/makerdao/vulcanizedb/pkg/core"

	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
)

var rawVatHealLog = types.Log{
	Address: common.HexToAddress(VatAddress()),
	Topics: []common.Hash{
		common.HexToHash("0x00000000000000000000000000000000000000000000000000000000f37ac61c"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000002711"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e0f37ac61c00000000000000000000000000000000000000000000000000000000000027110000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 105,
	TxHash:      common.HexToHash("0xe5bebbe1ea46d8b6d1515ce9ac4659f9e6774669c1f2009dcc7289c18b91e393"),
	TxIndex:     2,
	BlockHash:   common.HexToHash("0x3e011e52723db56476dc8cd45e1325f7bf3f3b2d89651253d6e8b66489f37d7c"),
	Index:       3,
	Removed:     false,
}

var VatHealHeaderSyncLog = core.HeaderSyncLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawVatHealLog,
	Transformed: false,
}

var VatHealModel = shared.InsertionModel{
	SchemaName: "maker",
	TableName:  "vat_heal",
	OrderedColumns: []string{
		constants.HeaderFK, "rad", constants.LogFK,
	},
	ColumnValues: shared.ColumnValues{
		"rad":              "10001",
		constants.HeaderFK: VatHealHeaderSyncLog.HeaderID,
		constants.LogFK:    VatHealHeaderSyncLog.ID,
	},
	ForeignKeyValues: shared.ForeignKeyValues{},
}
