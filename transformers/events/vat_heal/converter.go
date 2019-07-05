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

package vat_heal

import (
	"encoding/json"
	"errors"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/mcd_transformers/transformers/shared"
)

type VatHealConverter struct{}

func (VatHealConverter) ToModels(ethLogs []types.Log) ([]shared.InsertionModel, error) {
	var models []shared.InsertionModel
	for _, ethLog := range ethLogs {
		err := verifyLog(ethLog)
		if err != nil {
			return nil, err
		}

		radInt := shared.ConvertUint256HexToBigInt(ethLog.Topics[1].Hex())

		rawLogJson, err := json.Marshal(ethLog)
		if err != nil {
			return nil, err
		}

		model := shared.InsertionModel{
			TableName: "vat_heal",
			OrderedColumns: []string{
				"header_id", "rad", "log_idx", "tx_idx", "raw_log",
			},
			ColumnValues: shared.ColumnValues{
				"rad":     radInt.String(),
				"log_idx": ethLog.Index,
				"tx_idx":  ethLog.TxIndex,
				"raw_log": rawLogJson,
			},
			ForeignKeyValues: shared.ForeignKeyValues{},
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
