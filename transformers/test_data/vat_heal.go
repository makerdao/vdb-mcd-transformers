// VulcanizeDB
// Copyright © 2018 Vulcanize

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

	"github.com/vulcanize/mcd_transformers/transformers/events/vat_heal"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

var EthVatHealLogWithPositiveRad = types.Log{
	Address: common.HexToAddress(constants.VatContractAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.VatHealSignature()),
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

var rawVatHealLogWithPositiveRad, _ = json.Marshal(EthVatHealLogWithPositiveRad)
var VatHealModelWithPositiveRad = vat_heal.VatHealModel{
	Urn:              "0x07Fa9eF6609cA7921112231F8f195138ebbA2977",
	V:                "0x7340e006f4135BA6970D43bf43d88DCAD4e7a8CA",
	Rad:              "1000000000000000000000",
	LogIndex:         EthVatHealLogWithPositiveRad.Index,
	TransactionIndex: EthVatHealLogWithPositiveRad.TxIndex,
	Raw:              rawVatHealLogWithPositiveRad,
}

var EthVatHealLogWithNegativeRad = types.Log{
	Address: common.HexToAddress(constants.VatContractAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.VatHealSignature()),
		common.HexToHash("0x00000000000000000000000007fa9ef6609ca7921112231f8f195138ebba2977"),
		common.HexToHash("0x0000000000000000000000007340e006f4135ba6970d43bf43d88dcad4e7a8ca"),
		common.HexToHash("0xffffffffffffffffffffffffffffffffffffffffffffffc9ca36523a21600000"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e0ee8cd74800000000000000000000000007fa9ef6609ca7921112231f8f195138ebba29770000000000000000000000007340e006f4135ba6970d43bf43d88dcad4e7a8caffffffffffffffffffffffffffffffffffffffffffffffc9ca36523a2160000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 106,
	TxHash:      common.HexToHash("0xf9956cfa5f5290b087a99d9e667cd0b61ea80020901cd570293f6fac745b6eac"),
	TxIndex:     2,
	BlockHash:   common.HexToHash("0xeca99c865e695c1a8e2fae5df0b359da145f72eae8b4873d2da3e190213d0cf2"),
	Index:       3,
	Removed:     false,
}

var rawVatHealLogWithNegativeRad, _ = json.Marshal(EthVatHealLogWithNegativeRad)
var VatHealModelWithNegativeRad = vat_heal.VatHealModel{
	Urn:              "0x07Fa9eF6609cA7921112231F8f195138ebbA2977",
	V:                "0x7340e006f4135BA6970D43bf43d88DCAD4e7a8CA",
	Rad:              "-1000000000000000000000",
	LogIndex:         EthVatHealLogWithNegativeRad.Index,
	TransactionIndex: EthVatHealLogWithNegativeRad.TxIndex,
	Raw:              rawVatHealLogWithNegativeRad,
}
