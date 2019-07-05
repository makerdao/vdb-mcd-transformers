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

package dent

import (
	"encoding/json"
	"errors"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/mcd_transformers/transformers/shared"
)

type DentConverter struct{}

func (c DentConverter) ToModels(ethLogs []types.Log) (result []shared.InsertionModel, err error) {
	for _, log := range ethLogs {
		validateErr := validateLog(log)
		if validateErr != nil {
			return nil, validateErr
		}

		bidId := log.Topics[2].Big()
		lot := log.Topics[3].Big()
		bidBytes, dataErr := shared.GetLogNoteArgumentAtIndex(2, log.Data)
		if dataErr != nil {
			return nil, dataErr
		}
		bid := shared.ConvertUint256HexToBigInt(hexutil.Encode(bidBytes))

		logIndex := log.Index
		transactionIndex := log.TxIndex

		raw, err := json.Marshal(log)
		if err != nil {
			return nil, err
		}

		model := shared.InsertionModel{
			TableName: "dent",
			OrderedColumns: []string{
				"header_id", "bid_id", "lot", "bid", "contract_address", "tic", "log_idx", "tx_idx", "raw_log",
			},
			ColumnValues: shared.ColumnValues{
				"bid_id":           bidId.String(),
				"lot":              lot.String(),
				"bid":              bid.String(),
				"contract_address": log.Address.Hex(),
				"log_idx":          logIndex,
				"tx_idx":           transactionIndex,
				"raw_log":          raw,
			},
			ForeignKeyValues: shared.ForeignKeyValues{},
		}
		result = append(result, model)
	}
	return result, err
}

func validateLog(ethLog types.Log) error {
	if len(ethLog.Data) <= 0 {
		return errors.New("dent log data is empty")
	}

	if len(ethLog.Topics) < 4 {
		return errors.New("dent log does not contain expected topics")
	}

	return nil
}
