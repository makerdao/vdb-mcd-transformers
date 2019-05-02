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
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/transformers/events/pip_log_value"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

var (
	blockNumber = uint64(10606964)
	txIndex     = uint(2)
)

// https://kovan.etherscan.io/tx/0xbf700fefd1817c91c6a3e2dfa9f2f84f1c4d6a42c13f91ac6aa64bfc63f2e568
var EthPipLogValueLog = types.Log{
	Address:     common.HexToAddress(constants.PipEthContractAddress()),
	Topics:      []common.Hash{common.HexToHash(KovanPipLogValueSignature)},
	Data:        common.FromHex("0000000000000000000000000000000000000000000000076eec1981d1900000"),
	BlockNumber: blockNumber,
	TxHash:      common.HexToHash("0xbf700fefd1817c91c6a3e2dfa9f2f84f1c4d6a42c13f91ac6aa64bfc63f2e568"),
	TxIndex:     txIndex,
	BlockHash:   fakes.FakeHash,
	Index:       8,
	Removed:     false,
}

var rawPipLogValueLog, _ = json.Marshal(EthPipLogValueLog)
var PipLogValueModel = pip_log_value.PipLogValueModel{
	BlockNumber:      blockNumber,
	ContractAddress:  EthPipLogValueLog.Address.String(),
	Value:            "137120000000000000000",
	LogIndex:         EthPipLogValueLog.Index,
	TransactionIndex: EthPipLogValueLog.TxIndex,
	Raw:              rawPipLogValueLog,
}

func GetFakePipLogValue(blockNum int64, txIdx int, value string) pip_log_value.PipLogValueModel {
	var EthPipLogValueLog = types.Log{
		Address:     pipAddress,
		Topics:      []common.Hash{common.HexToHash(KovanPipLogValueSignature)},
		Data:        common.FromHex("0000000000000000000000000000000000000000000000076eec1981d1900000"),
		BlockNumber: uint64(blockNum),
		TxHash:      common.HexToHash("0xbf700fefd1817c91c6a3e2dfa9f2f84f1c4d6a42c13f91ac6aa64bfc63f2e568"),
		TxIndex:     uint(txIdx),
		BlockHash:   fakes.FakeHash,
		Index:       8,
		Removed:     false,
	}

	var rawPipLogValue, _ = json.Marshal(EthPipLogValueLog)

	return pip_log_value.PipLogValueModel{
		BlockNumber:      uint64(blockNum),
		ContractAddress:  EthPipLogValueLog.Address.String(),
		Value:            value,
		LogIndex:         EthPipLogValueLog.Index,
		TransactionIndex: EthPipLogValueLog.TxIndex,
		Raw:              rawPipLogValue,
	}
}
