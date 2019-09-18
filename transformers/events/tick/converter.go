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

package tick

import (
	"encoding/json"
	"errors"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/mcd_transformers/transformers/shared"
)

type TickConverter struct{}

func (TickConverter) ToModels(ethLogs []types.Log) (results []shared.InsertionModel, err error) {
	for _, ethLog := range ethLogs {
		validateErr := validateLog(ethLog)
		if validateErr != nil {
			return nil, validateErr
		}

		rawLog, jsonErr := json.Marshal(ethLog)
		if jsonErr != nil {
			return nil, jsonErr
		}

		model := shared.InsertionModel{
			SchemaName: "maker",
			TableName:  "tick",
			OrderedColumns: []string{
				"header_id", "bid_id", string(constants.AddressFK), "log_idx", "tx_idx", "raw_log",
			},
			ColumnValues: shared.ColumnValues{
				"bid_id":  ethLog.Topics[2].Big().String(),
				"log_idx": ethLog.Index,
				"tx_idx":  ethLog.TxIndex,
				"raw_log": rawLog,
			},
			ForeignKeyValues: shared.ForeignKeyValues{
				constants.AddressFK: ethLog.Address.String(),
			},
		}
		results = append(results, model)
	}
	return results, err
}

func validateLog(ethLog types.Log) error {
	if len(ethLog.Topics) < 3 {
		return errors.New("flip tick log does not contain expected topics")
	}

	return nil
}
