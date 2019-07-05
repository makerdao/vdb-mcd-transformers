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

package tend

import (
	"encoding/json"
	"errors"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/mcd_transformers/transformers/shared"
)

type TendConverter struct{}

func (TendConverter) ToModels(ethLogs []types.Log) (results []shared.InsertionModel, err error) {
	for _, ethLog := range ethLogs {
		err := validateLog(ethLog)
		if err != nil {
			return nil, err
		}

		bidId := ethLog.Topics[2].Big()
		lot := ethLog.Topics[3].Big().String()
		rawBid, bidErr := shared.GetLogNoteArgumentAtIndex(2, ethLog.Data)
		if bidErr != nil {
			return nil, bidErr
		}
		bidValue := shared.ConvertUint256HexToBigInt(hexutil.Encode(rawBid)).String()
		contractAddress := ethLog.Address
		logIndex := ethLog.Index
		transactionIndex := ethLog.TxIndex

		rawLog, err := json.Marshal(ethLog)
		if err != nil {
			return nil, err
		}

		model := shared.InsertionModel{
			TableName: "tend",
			OrderedColumns: []string{
				"header_id", "bid_id", "lot", "bid", "contract_address", "log_idx", "tx_idx", "raw_log",
			},
			ColumnValues: shared.ColumnValues{
				"bid_id":           bidId.String(),
				"lot":              lot,
				"bid":              bidValue,
				"contract_address": contractAddress.Hex(),
				"log_idx":          logIndex,
				"tx_idx":           transactionIndex,
				"raw_log":          rawLog,
			},
			ForeignKeyValues: shared.ForeignKeyValues{},
		}
		results = append(results, model)
	}
	return results, err
}

func validateLog(ethLog types.Log) error {
	if len(ethLog.Data) <= 0 {
		return errors.New("tend log note data is empty")
	}

	if len(ethLog.Topics) < 4 {
		return errors.New("tend log does not contain expected topics")
	}

	return nil
}
