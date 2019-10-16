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

var rawSpotFileMatLog = types.Log{
	Address: common.HexToAddress(SpotAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.SpotFileMatSignature()),
		common.HexToHash("0x00000000000000000000000071ce79fcfec71760d51f6b3589c0d9ec0ccb64a8"),
		common.HexToHash("0x4554482d41000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x6d61740000000000000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e01a0b287e4554482d410000000000000000000000000000000000000000000000000000006d61740000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004d8c55aefb8c05b5c00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 11257385,
	TxHash:      common.HexToHash("0xb4d19aaf5be5077db65aeeb16906a5498cfa94836952191258cc45966e1d7198"),
	TxIndex:     2,
	BlockHash:   common.HexToHash("0x968cd16acb356de42e9f3ab17583988b49173c0339af5afa3f79cecdbc111d69"),
	Index:       3,
	Removed:     false,
}

var SpotFileMatHeaderSyncLog = core.HeaderSyncLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawSpotFileMatLog,
	Transformed: false,
}

func SpotFileMatModel() shared.InsertionModel { return CopyModel(spotFileMatModel) }

var spotFileMatModel = shared.InsertionModel{
	SchemaName: "maker",
	TableName:  "spot_file_mat",
	OrderedColumns: []string{
		constants.HeaderFK, string(constants.IlkFK), "what", "data", constants.LogFK,
	},
	ColumnValues: shared.ColumnValues{
		"what":             "mat",
		"data":             "1500000000000000000000000000",
		constants.HeaderFK: SpotFileMatHeaderSyncLog.HeaderID,
		constants.LogFK:    SpotFileMatHeaderSyncLog.ID,
	},
	ForeignKeyValues: shared.ForeignKeyValues{
		constants.IlkFK: "0x4554482d41000000000000000000000000000000000000000000000000000000",
	},
}

var rawSpotFilePipLog = types.Log{
	Address: common.HexToAddress(SpotAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.SpotFilePipSignature()),
		common.HexToHash("0x00000000000000000000000073acbfb5b9413b0020164ee63dce4e1f71aba67c"),
		common.HexToHash("0x4554482d41000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x7069700000000000000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e0ebecb39d4554482d41000000000000000000000000000000000000000000000000000000706970000000000000000000000000000000000000000000000000000000000000000000000000000000000075dd74e8afe8110c8320ed397cccff3b8134d98100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 11257191,
	TxHash:      common.HexToHash("0xaae9e8bce346f86a01c5a3af137bc1f9bc7c0c767804a2b9b6356849aee0d7dd"),
	TxIndex:     1,
	BlockHash:   common.HexToHash("0xfa28e186578238fdd6b971add2ebe62a26dddf5ff971d50ee476c86b45362da1"),
	Index:       2,
	Removed:     false,
}

var SpotFilePipHeaderSyncLog = core.HeaderSyncLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawSpotFilePipLog,
	Transformed: false,
}

func SpotFilePipModel() shared.InsertionModel { return CopyModel(spotFilePipModel) }

var spotFilePipModel = shared.InsertionModel{
	SchemaName: "maker",
	TableName:  "spot_file_pip",
	OrderedColumns: []string{
		constants.HeaderFK, string(constants.IlkFK), "what", "pip", constants.LogFK,
	},
	ColumnValues: shared.ColumnValues{
		"what":             "pip",
		"pip":              "0x75dD74e8afE8110C8320eD397CcCff3B8134d981",
		constants.HeaderFK: SpotFilePipHeaderSyncLog.HeaderID,
		constants.LogFK:    SpotFilePipHeaderSyncLog.ID,
	},
	ForeignKeyValues: shared.ForeignKeyValues{
		constants.IlkFK: "0x4554482d41000000000000000000000000000000000000000000000000000000",
	},
}
