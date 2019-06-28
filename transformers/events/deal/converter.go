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

package deal

import (
	"encoding/json"
	"errors"

	"github.com/ethereum/go-ethereum/core/types"
)

type DealConverter struct{}

func (DealConverter) ToModels(ethLogs []types.Log) (result []interface{}, err error) {
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

		model := DealModel{
			BidId:            bidId.String(),
			ContractAddress:  log.Address.Hex(),
			LogIndex:         log.Index,
			TransactionIndex: log.TxIndex,
			Raw:              raw,
		}
		result = append(result, model)
	}

	return result, nil
}

func validateLog(ethLog types.Log) error {
	if len(ethLog.Topics) < 3 {
		return errors.New("deal log does not contain expected topics")
	}
	return nil
}
