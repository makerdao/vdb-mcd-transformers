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

package vat_move

import (
	"encoding/json"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/mcd_transformers/transformers/shared"
)

type VatMoveConverter struct{}

func (VatMoveConverter) ToModels(ethLogs []types.Log) ([]shared.InsertionModel, error) {
	var models []shared.InsertionModel
	for _, ethLog := range ethLogs {
		err := verifyLog(ethLog)
		if err != nil {
			return []shared.InsertionModel{}, err
		}

		src := common.BytesToAddress(ethLog.Topics[1].Bytes()).String()
		dst := common.BytesToAddress(ethLog.Topics[2].Bytes()).String()
		rad := shared.ConvertUint256HexToBigInt(ethLog.Topics[3].Hex())
		raw, err := json.Marshal(ethLog)
		if err != nil {
			return []shared.InsertionModel{}, err
		}

		model := shared.InsertionModel{
			TableName: "vat_move",
			OrderedColumns: []string{
				"header_id", "src", "dst", "rad", "log_idx", "tx_idx", "raw_log",
			},
			ColumnValues: shared.ColumnValues{
				"src":     src,
				"dst":     dst,
				"rad":     rad.String(),
				"log_idx": ethLog.Index,
				"tx_idx":  ethLog.TxIndex,
				"raw_log": raw,
			},
			ForeignKeyValues: shared.ForeignKeyValues{},
		}
		models = append(models, model)
	}

	return models, nil
}

func verifyLog(ethLog types.Log) error {
	if len(ethLog.Data) <= 0 {
		return errors.New("log data is empty")
	}
	if len(ethLog.Topics) < 4 {
		return errors.New("log missing topics")
	}
	return nil
}
