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

var EthVatGrabLogWithPositiveDink = types.Log{
	Address: common.HexToAddress(constants.VatContractAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.VatGrabSignature()),
		common.HexToHash("0x66616b6520696c6b000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x00000000000000000000000007fa9ef6609ca7921112231f8f195138ebba2977"),
		common.HexToHash("0x0000000000000000000000007340e006f4135ba6970d43bf43d88dcad4e7a8ca"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e07bab3f4066616b6520696c6b00000000000000000000000000000000000000000000000000000000000000000000000007fa9ef6609ca7921112231f8f195138ebba29770000000000000000000000007340e006f4135ba6970d43bf43d88dcad4e7a8ca0000000000000000000000007526eb4f95e2a1394797cb38a921fb1eba09291b00000000000000000000000000000000000000000000003635c9adc5dea0000000000000000000000000000000000000000000000000006c6b935b8bbd40000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 94,
	TxHash:      common.HexToHash("0x18aa10dddd3543143d4eff314c8f3620c3b8299e57468d70ff7abd4498eee7db"),
	TxIndex:     4,
	BlockHash:   common.HexToHash("0x17b1de2797689d940a66911ebb9ae789528c269aea309c55128e38d48ddb37a9"),
	Index:       5,
	Removed:     false,
}

var rawVatGrabLogWithPositiveDink, _ = json.Marshal(EthVatGrabLogWithPositiveDink)
var VatGrabModelWithPositiveDink = shared.InsertionModel{
	TableName: "vat_grab",
	OrderedColumns: []string{
		"header_id", string(constants.UrnFK), "v", "w", "dink", "dart", "log_idx", "tx_idx", "raw_log",
	},
	ColumnValues: shared.ColumnValues{
		"v":       "0x7340e006f4135BA6970D43bf43d88DCAD4e7a8CA",
		"w":       "0x7526EB4f95e2a1394797Cb38a921Fb1EbA09291B",
		"dink":    "1000000000000000000000",
		"dart":    "2000000000000000000000",
		"log_idx": EthVatGrabLogWithPositiveDink.Index,
		"tx_idx":  EthVatGrabLogWithPositiveDink.TxIndex,
		"raw_log": rawVatGrabLogWithPositiveDink,
	},
	ForeignKeyValues: shared.ForeignKeyValues{
		constants.IlkFK: "0x66616b6520696c6b000000000000000000000000000000000000000000000000",
		constants.UrnFK: "0x07Fa9eF6609cA7921112231F8f195138ebbA2977",
	},
}

var EthVatGrabLogWithNegativeDink = types.Log{
	Address: common.HexToAddress(constants.VatContractAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.VatGrabSignature()),
		common.HexToHash("0x66616b6520696c6b000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x00000000000000000000000007fa9ef6609ca7921112231f8f195138ebba2977"),
		common.HexToHash("0x0000000000000000000000007340e006f4135ba6970d43bf43d88dcad4e7a8ca"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e07bab3f4066616b6520696c6b00000000000000000000000000000000000000000000000000000000000000000000000007fa9ef6609ca7921112231f8f195138ebba29770000000000000000000000007340e006f4135ba6970d43bf43d88dcad4e7a8ca0000000000000000000000007526eb4f95e2a1394797cb38a921fb1eba09291bffffffffffffffffffffffffffffffffffffffffffffffc9ca36523a21600000ffffffffffffffffffffffffffffffffffffffffffffff93946ca47442c0000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 95,
	TxHash:      common.HexToHash("0x38c279812f842d7498ffe2f45f0d8c5de76ecd1ee7912636e7431129746b2e15"),
	TxIndex:     0,
	BlockHash:   common.HexToHash("0x4ea9f55ed4a97f686dd28a0d6d2dea8a0915a902a245e86516df3e4a57d1ca9d"),
	Index:       0,
	Removed:     false,
}

var rawVatGrabLogWithNegativeDink, _ = json.Marshal(EthVatGrabLogWithNegativeDink)
var VatGrabModelWithNegativeDink = shared.InsertionModel{
	TableName: "vat_grab",
	OrderedColumns: []string{
		"header_id", string(constants.UrnFK), "v", "w", "dink", "dart", "log_idx", "tx_idx", "raw_log",
	},
	ColumnValues: shared.ColumnValues{
		"v":       "0x7340e006f4135BA6970D43bf43d88DCAD4e7a8CA",
		"w":       "0x7526EB4f95e2a1394797Cb38a921Fb1EbA09291B",
		"dink":    "-1000000000000000000000",
		"dart":    "-2000000000000000000000",
		"log_idx": EthVatGrabLogWithNegativeDink.Index,
		"tx_idx":  EthVatGrabLogWithNegativeDink.TxIndex,
		"raw_log": rawVatGrabLogWithNegativeDink,
	},
	ForeignKeyValues: shared.ForeignKeyValues{
		constants.IlkFK: "0x66616b6520696c6b000000000000000000000000000000000000000000000000",
		constants.UrnFK: "0x07Fa9eF6609cA7921112231F8f195138ebbA2977",
	},
}
