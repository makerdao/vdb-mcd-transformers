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
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_frob"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

var (
	frobData = "0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000c441fd3ef74554480000000000000000000000000000000000000000000000000000000000da15dce70ab462e66779f23ee14f21d993789ee3000000000000000000000000da15dce70ab462e66779f23ee14f21d993789ee3000000000000000000000000da15dce70ab462e66779f23ee14f21d993789ee30000000000000000000000000000000000000000000000000000000000000000000000000de0b6b3a7640000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
	frobUrn  = "da15dce70ab462e66779f23ee14f21d993789ee3000000000000000000000000"
)

var EthVatFrobLog = types.Log{
	Address: common.HexToAddress(KovanVatContractAddress),
	Topics: []common.Hash{
		common.HexToHash(KovanVatFrobSignature),
		common.HexToHash("0x4554480000000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0xda15dce70ab462e66779f23ee14f21d993789ee3000000000000000000000000"),
		common.HexToHash("0xda15dce70ab462e66779f23ee14f21d993789ee3000000000000000000000000"),
	},
	Data:        hexutil.MustDecode(frobData),
	BlockNumber: 10512592,
	TxHash:      common.HexToHash("0x10277b770bcd569cd3c943db2228153435ee1320eaab1f3a64fb8d5732d44c2e"),
	TxIndex:     123,
	BlockHash:   fakes.FakeHash,
	Index:       7,
	Removed:     false,
}

var rawVatFrobLog, _ = json.Marshal(EthVatFrobLog)
var VatFrobModel = vat_frob.VatFrobModel{
	Ilk:              "4554480000000000000000000000000000000000000000000000000000000000",
	Urn:              frobUrn,
	V:                frobUrn,
	W:                frobUrn,
	Dink:             "1000000000000000000",
	Dart:             "0",
	LogIndex:         EthVatFrobLog.Index,
	TransactionIndex: EthVatFrobLog.TxIndex,
	Raw:              rawVatFrobLog,
}
