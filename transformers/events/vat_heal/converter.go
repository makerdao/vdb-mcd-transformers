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

package vat_heal

import (
	"encoding/json"
	"errors"
	"github.com/vulcanize/mcd_transformers/transformers/shared"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/common"
)

type VatHealConverter struct{}

func (VatHealConverter) ToModels(ethLogs []types.Log) ([]interface{}, error) {
	var models []interface{}
	for _, ethLog := range ethLogs {
		err := verifyLog(ethLog)
		if err != nil {
			return nil, err
		}

		urn := common.BytesToAddress(ethLog.Topics[1].Bytes()).String()
		v := common.BytesToAddress(ethLog.Topics[2].Bytes()).String()
		radInt := shared.ConvertInt256HexToBigInt(ethLog.Topics[3].Hex())

		rawLogJson, err := json.Marshal(ethLog)
		if err != nil {
			return nil, err
		}

		model := VatHealModel{
			Urn:              urn,
			V:                v,
			Rad:              radInt.String(),
			LogIndex:         ethLog.Index,
			TransactionIndex: ethLog.TxIndex,
			Raw:              rawLogJson,
		}

		models = append(models, model)
	}

	return models, nil
}

func verifyLog(log types.Log) error {
	if len(log.Topics) < 4 {
		return errors.New("log missing topics")
	}
	return nil
}
