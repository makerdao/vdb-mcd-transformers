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

package pip_log_value

import (
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
)

type PipLogValueConverter struct{}

func (converter PipLogValueConverter) ToModels(ethLogs []types.Log) ([]interface{}, error) {
	var results []interface{}
	for _, log := range ethLogs {
		raw, err := json.Marshal(log)
		if err != nil {
			return nil, err
		}
		value := new(big.Int).SetBytes(log.Data)
		model := PipLogValueModel{
			BlockNumber:      log.BlockNumber,
			ContractAddress:  log.Address.String(),
			Value:            value.String(),
			LogIndex:         log.Index,
			TransactionIndex: log.TxIndex,
			Raw:              raw,
		}
		results = append(results, model)
	}
	return results, nil
}
