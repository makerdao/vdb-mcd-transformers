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
	"github.com/vulcanize/mcd_transformers/transformers/events/cat_file/chop_lump"
	"github.com/vulcanize/mcd_transformers/transformers/events/cat_file/flip"
	"github.com/vulcanize/mcd_transformers/transformers/events/cat_file/vow"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

var EthCatFileChopLog = types.Log{
	Address: common.HexToAddress(KovanCatContractAddress),
	Topics: []common.Hash{
		common.HexToHash(KovanCatFileChopLumpSignature),
		common.HexToHash("0x00000000000000000000000064d922894153be9eef7b7218dc565d1d0ce2a092"),
		common.HexToHash("0x66616b6520696c6b000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x63686f7000000000000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000641a0b287e524550000000000000000000000000000000000000000000000000000000000063686f70000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000033b2e3c9fd0803ce800000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 110,
	TxHash:      common.HexToHash("0xe32dfe6afd7ea28475569756fc30f0eea6ad4cfd32f67436ff1d1c805e4382df"),
	TxIndex:     13,
	BlockHash:   fakes.FakeHash,
	Index:       1,
	Removed:     false,
}
var rawCatFileChopLog, _ = json.Marshal(EthCatFileChopLog)
var CatFileChopModel = chop_lump.CatFileChopLumpModel{
	Ilk:              "66616b6520696c6b000000000000000000000000000000000000000000000000",
	What:             "chop",
	Data:             "1000000000000000000000000000",
	TransactionIndex: EthCatFileChopLog.TxIndex,
	LogIndex:         EthCatFileChopLog.Index,
	Raw:              rawCatFileChopLog,
}

var EthCatFileLumpLog = types.Log{
	Address: common.HexToAddress(KovanCatContractAddress),
	Topics: []common.Hash{
		common.HexToHash(KovanCatFileChopLumpSignature),
		common.HexToHash("0x00000000000000000000000064d922894153be9eef7b7218dc565d1d0ce2a092"),
		common.HexToHash("0x66616b6520696c6b000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x6c756d7000000000000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000641a0b287e52455000000000000000000000000000000000000000000000000000000000006c756d700000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000021e19e0c9bab240000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 110,
	TxHash:      common.HexToHash("0xe32dfe6afd7ea28475569756fc30f0eea6ad4cfd32f67436ff1d1c805e4382df"),
	TxIndex:     15,
	BlockHash:   fakes.FakeHash,
	Index:       3,
	Removed:     false,
}
var rawCatFileLumpLog, _ = json.Marshal(EthCatFileLumpLog)
var CatFileLumpModel = chop_lump.CatFileChopLumpModel{
	Ilk:              "66616b6520696c6b000000000000000000000000000000000000000000000000",
	What:             "lump",
	Data:             "10000000000000000000000",
	TransactionIndex: EthCatFileLumpLog.TxIndex,
	LogIndex:         EthCatFileLumpLog.Index,
	Raw:              rawCatFileLumpLog,
}

var EthCatFileFlipLog = types.Log{
	Address: common.HexToAddress(KovanCatContractAddress),
	Topics: []common.Hash{
		common.HexToHash(KovanCatFileFlipSignature),
		common.HexToHash("0x00000000000000000000000064d922894153be9eef7b7218dc565d1d0ce2a092"),
		common.HexToHash("0x66616b6520696c6b000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x666c697000000000000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000064ebecb39d5245500000000000000000000000000000000000000000000000000000000000666c6970000000000000000000000000000000000000000000000000000000000000000000000000000000004ec982bc57c463d4a1825d975e2a525c4daadd9100000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 88,
	TxHash:      common.HexToHash("0xc71ef3e9999595913d31e89446cab35319bd4289520e55611a1b42fc2a8463b6"),
	TxIndex:     12,
	BlockHash:   fakes.FakeHash,
	Index:       1,
	Removed:     false,
}

var rawCatFileFlipLog, _ = json.Marshal(EthCatFileFlipLog)
var CatFileFlipModel = flip.CatFileFlipModel{
	Ilk:              "66616b6520696c6b000000000000000000000000000000000000000000000000",
	What:             "flip",
	Flip:             "0x4EC982bC57c463D4A1825d975E2A525C4daadD91",
	TransactionIndex: EthCatFileFlipLog.TxIndex,
	LogIndex:         EthCatFileFlipLog.Index,
	Raw:              rawCatFileFlipLog,
}

var EthCatFileVowLog = types.Log{
	Address: common.HexToAddress(KovanCatContractAddress),
	Topics: []common.Hash{
		common.HexToHash(KovanCatFileVowSignature),
		common.HexToHash("0x0000000000000000000000003652c2af10cbbdb753c3b46489db5226b73e6497"),
		common.HexToHash("0x766f770000000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x00000000000000000000000017560834075da3db54f737db74377e799c865821"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000044d4e8be83766f77000000000000000000000000000000000000000000000000000000000000000000000000000000000017560834075da3db54f737db74377e799c86582100000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 87,
	TxHash:      common.HexToHash("0x6515c7dfe53f0ad83ce1173fa99032c24a07cfd8b5d5a1c1f80486c99dd52800"),
	TxIndex:     11,
	BlockHash:   fakes.FakeHash,
	Index:       2,
	Removed:     false,
}

var rawCatFileVowLog, _ = json.Marshal(EthCatFileVowLog)
var CatFileVowModel = vow.CatFileVowModel{
	What:             "vow",
	Data:             "0x17560834075DA3Db54f737db74377E799c865821",
	TransactionIndex: EthCatFileVowLog.TxIndex,
	LogIndex:         EthCatFileVowLog.Index,
	Raw:              rawCatFileVowLog,
}
