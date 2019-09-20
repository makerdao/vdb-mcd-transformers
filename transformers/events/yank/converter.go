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

package yank

import (
	"encoding/json"
	"errors"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
)

type YankConverter struct{}

func (YankConverter) ToModels(_ string, ethLogs []types.Log) (results []shared.InsertionModel, err error) {
	for _, log := range ethLogs {
		validationErr := validateLog(log)
		if validationErr != nil {
			return nil, validationErr
		}

		bidId := log.Topics[2].Big()
		raw, jsonErr := json.Marshal(log)
		if jsonErr != nil {
			return nil, jsonErr
		}

		model := shared.InsertionModel{
			SchemaName: "maker",
			TableName:  "yank",
			OrderedColumns: []string{
				"header_id", "bid_id", string(constants.AddressFK), "log_idx", "tx_idx", "raw_log",
			},
			ColumnValues: shared.ColumnValues{
				"bid_id":  bidId.String(),
				"log_idx": log.Index,
				"tx_idx":  log.TxIndex,
				"raw_log": raw,
			},
			ForeignKeyValues: shared.ForeignKeyValues{
				constants.AddressFK: log.Address.Hex(),
			},
		}
		results = append(results, model)
	}
	return results, err
}

func validateLog(ethLog types.Log) error {
	if len(ethLog.Topics) < 3 {
		return errors.New("yank log does not contain expected topics")
	}
	return nil
}
