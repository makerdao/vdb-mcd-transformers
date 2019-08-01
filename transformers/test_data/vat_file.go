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
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

var EthVatFileDebtCeilingLog = types.Log{
	Address: common.HexToAddress(VatAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.VatFileDebtCeilingSignature()),
		common.HexToHash("0x4c696e6500000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x0000000000000000000002ac3a4edbbfb8014e3ba83411e915e8000000000000"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000004429ae81144c696e65000000000000000000000000000000000000000000000000000000000000000000000000000002ac3a4edbbfb8014e3ba83411e915e800000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 10577169,
	TxHash:      common.HexToHash("0x0ec18121d45f96293d9d759fd7564db4186a1aa69552f3106dd1afdeffdc9106"),
	TxIndex:     333,
	BlockHash:   fakes.FakeHash,
	Index:       15,
	Removed:     false,
}

var rawVatFileDebtCeilingLog, _ = json.Marshal(EthVatFileDebtCeilingLog)
var VatFileDebtCeilingModel = shared.InsertionModel{
	TableName: "vat_file_debt_ceiling",
	OrderedColumns: []string{
		"header_id", "what", "data", "log_idx", "tx_idx", "raw_log",
	},
	ColumnValues: shared.ColumnValues{
		"what":    "Line",
		"data":    "1000000000000000000000000000000000000000000000000000",
		"log_idx": EthVatFileDebtCeilingLog.Index,
		"tx_idx":  EthVatFileDebtCeilingLog.TxIndex,
		"raw_log": rawVatFileDebtCeilingLog,
	},
	ForeignKeyValues: shared.ForeignKeyValues{},
}

var EthVatFileIlkDustLog = types.Log{
	Address: common.HexToAddress(VatAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.VatFileIlkSignature()),                                       //sig
		common.HexToHash("0x5245500000000000000000000000000000000000000000000000000000000000"),  //ilk
		common.HexToHash("0x6475737400000000000000000000000000000000000000000000000000000000"),  //what
		common.HexToHash("0x000000000000000000000000000000000000000832600000bee4f14727b555555"), //val
	},

	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000641a0b287e52455000000000000000000000000000000000000000000000000000000000006475737400000000000000000000000000000000000000000000000000000000000000000000000000000000000000832600000bee4f14727b55555500000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 12,
	TxHash:      common.HexToHash("0x2e27c962a697d4f7ec5d3206d0c058bd510f7593a711f082e55f3b62d44d8dd9"),
	TxIndex:     112,
	BlockHash:   fakes.FakeHash,
	Index:       15,
	Removed:     false,
}

var rawVatFileIlkDustLog, _ = json.Marshal(EthVatFileIlkDustLog)
var VatFileIlkDustModel = shared.InsertionModel{
	TableName: "vat_file_ilk",
	OrderedColumns: []string{
		"header_id", string(constants.IlkFK), "what", "data", "log_idx", "tx_idx", "raw_log",
	},
	ColumnValues: shared.ColumnValues{
		"what":    "dust",
		"data":    "10390649719961925488562719249749",
		"log_idx": EthVatFileIlkDustLog.Index,
		"tx_idx":  EthVatFileIlkDustLog.TxIndex,
		"raw_log": rawVatFileIlkDustLog,
	},
	ForeignKeyValues: shared.ForeignKeyValues{
		constants.IlkFK: "0x5245500000000000000000000000000000000000000000000000000000000000",
	},
}

var EthVatFileIlkLineLog = types.Log{
	Address: common.HexToAddress(VatAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.VatFileIlkSignature()),
		common.HexToHash("0x5245500000000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x6c696e6500000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x00000000000000000000000000000000000000000000000bee4f14727b555555"),
	},

	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000641a0b287e52455000000000000000000000000000000000000000000000000000000000006c696e6500000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000bee4f14727b55555500000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 12,
	TxHash:      common.HexToHash("0x2e27c962a697d4f7ec5d3206d0c058bd510f7593a711f082e55f3b62d44d8dd9"),
	TxIndex:     112,
	BlockHash:   fakes.FakeHash,
	Index:       15,
	Removed:     false,
}

var rawVatFileIlkLineLog, _ = json.Marshal(EthVatFileIlkLineLog)
var VatFileIlkLineModel = shared.InsertionModel{
	TableName: "vat_file_ilk",
	OrderedColumns: []string{
		"header_id", string(constants.IlkFK), "what", "data", "log_idx", "tx_idx", "raw_log",
	},
	ColumnValues: shared.ColumnValues{
		"what":    "line",
		"data":    "220086151196920075605",
		"log_idx": EthVatFileIlkLineLog.Index,
		"tx_idx":  EthVatFileIlkLineLog.TxIndex,
		"raw_log": rawVatFileIlkLineLog,
	},
	ForeignKeyValues: shared.ForeignKeyValues{
		constants.IlkFK: "0x5245500000000000000000000000000000000000000000000000000000000000",
	},
}

var EthVatFileIlkSpotLog = types.Log{
	Address: common.HexToAddress(VatAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.VatFileIlkSignature()),
		common.HexToHash("0x4554480000000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x73706f7400000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x00000000000000000000000000000000000000012714e40bee4f14727b555555"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000641a0b287e455448000000000000000000000000000000000000000000000000000000000073706f740000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000012714e40bee4f14727b55555500000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 10501158,
	TxHash:      common.HexToHash("0x657baea11882f6a73b0088382fa9b9b7ba84f0c1183af909d93ab0fe2d10c292"),
	TxIndex:     111,
	BlockHash:   fakes.FakeHash,
	Index:       14,
	Removed:     false,
}

var rawVatFileIlkSpotLog, _ = json.Marshal(EthVatFileIlkSpotLog)
var VatFileIlkSpotModel = shared.InsertionModel{
	TableName: "vat_file_ilk",
	OrderedColumns: []string{
		"header_id", string(constants.IlkFK), "what", "data", "log_idx", "tx_idx", "raw_log",
	},
	ColumnValues: shared.ColumnValues{
		"what":    "spot",
		"data":    "91323333333333333333333333333",
		"log_idx": EthVatFileIlkSpotLog.Index,
		"tx_idx":  EthVatFileIlkSpotLog.TxIndex,
		"raw_log": rawVatFileIlkSpotLog,
	},
	ForeignKeyValues: shared.ForeignKeyValues{
		constants.IlkFK: "0x4554480000000000000000000000000000000000000000000000000000000000",
	},
}
