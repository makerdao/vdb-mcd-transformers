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

	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

var EthVatForkLogWithNegativeDinkDart = types.Log{
	Address: common.HexToAddress(constants.GetContractAddress("MCD_VAT")),
	Topics: []common.Hash{
		common.HexToHash("0x00000000000000000000000000000000000000000000000000000000870c616d"),
		common.HexToHash("0x66616b6520696c6b000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x00000000000000000000000007fa9ef6609ca7921112231f8f195138ebba2977"),
		common.HexToHash("0x0000000000000000000000007526eb4f95e2a1394797cb38a921fb1eba09291b"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e0870c616d66616b6520696c6b00000000000000000000000000000000000000000000000000000000000000000000000007fa9ef6609ca7921112231f8f195138ebba29770000000000000000000000007526eb4f95e2a1394797cb38a921fb1eba09291bffffffffffffffffffffffffffffffffffffffffffffffc9ca36523a21600000ffffffffffffffffffffffffffffffffffffffffffffff93946ca47442c00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 146,
	TxHash:      common.HexToHash("0xe64bdd39a752e1911e841d634a6fa8d4ef229a03f0555f9e055caec1ae4930c2"),
	TxIndex:     1,
	BlockHash:   common.HexToHash("0xf31c6d2dadd23f408e5158dce47ba20fef8c17bc60af6e1f35a89769bc20d6f0"),
	Index:       2,
	Removed:     false,
}

var rawVatForkLogWithNegativeDinkDart, _ = json.Marshal(EthVatForkLogWithNegativeDinkDart)
var VatForkModelWithNegativeDinkDart = shared.InsertionModel{
	TableName: "vat_fork",
	OrderedColumns: []string{
		"header_id", string(constants.IlkFK), "src", "dst", "dink", "dart", "log_idx", "tx_idx", "raw_log",
	},
	ColumnValues: shared.ColumnValues{
		"src":     "0x07Fa9eF6609cA7921112231F8f195138ebbA2977",
		"dst":     "0x7526EB4f95e2a1394797Cb38a921Fb1EbA09291B",
		"dink":    "-1000000000000000000000",
		"dart":    "-2000000000000000000000",
		"log_idx": EthVatForkLogWithNegativeDinkDart.Index,
		"tx_idx":  EthVatForkLogWithNegativeDinkDart.TxIndex,
		"raw_log": rawVatForkLogWithNegativeDinkDart,
	},
	ForeignKeyValues: shared.ForeignKeyValues{
		constants.IlkFK: "0x66616b6520696c6b000000000000000000000000000000000000000000000000",
	},
}

var EthVatForkLogWithPositiveDinkDart = types.Log{
	Address: common.HexToAddress(constants.GetContractAddress("MCD_VAT")),
	Topics: []common.Hash{
		common.HexToHash("0x00000000000000000000000000000000000000000000000000000000870c616d"),
		common.HexToHash("0x66616b6520696c6b000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x00000000000000000000000007fa9ef6609ca7921112231f8f195138ebba2977"),
		common.HexToHash("0x0000000000000000000000007526eb4f95e2a1394797cb38a921fb1eba09291b"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e0870c616d66616b6520696c6b000000000000000000000000000000000000000000000000000000000000000000000000659344c807415e6d9f0d5b9cf61556d9d1dc4e8b0000000000000000000000000ccbc2b468e788e024553f105a30c84b1587e22500000000000000000000000000000000000000000000000000005af3107a400000000000000000000000000000000000000000000000000000071afd498d0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 146,
	TxHash:      common.HexToHash("0xe64bdd39a752e1911e841d634a6fa8d4ef229a03f0555f9e055caec1ae4930c2"),
	TxIndex:     1,
	BlockHash:   common.HexToHash("0xf31c6d2dadd23f408e5158dce47ba20fef8c17bc60af6e1f35a89769bc20d6f0"),
	Index:       2,
	Removed:     false,
}

var rawVatForkLogWithPositiveDinkDart, _ = json.Marshal(EthVatForkLogWithPositiveDinkDart)
var VatForkModelWithPositiveDinkDart = shared.InsertionModel{
	TableName: "vat_fork",
	OrderedColumns: []string{
		"header_id", string(constants.IlkFK), "src", "dst", "dink", "dart", "log_idx", "tx_idx", "raw_log",
	},
	ColumnValues: shared.ColumnValues{
		"src":     "0x07Fa9eF6609cA7921112231F8f195138ebbA2977",
		"dst":     "0x7526EB4f95e2a1394797Cb38a921Fb1EbA09291B",
		"dink":    "100000000000000",
		"dart":    "2000000000000000",
		"log_idx": EthVatForkLogWithPositiveDinkDart.Index,
		"tx_idx":  EthVatForkLogWithPositiveDinkDart.TxIndex,
		"raw_log": rawVatForkLogWithPositiveDinkDart,
	},
	ForeignKeyValues: shared.ForeignKeyValues{
		constants.IlkFK: "0x66616b6520696c6b000000000000000000000000000000000000000000000000",
	},
}
