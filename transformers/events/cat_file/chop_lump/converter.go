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

package chop_lump

import (
	"encoding/json"
	"errors"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/libraries/shared/constants"

	"github.com/vulcanize/mcd_transformers/transformers/shared"
	constants2 "github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

var (
	chop = "chop"
	lump = "lump"
)

type CatFileChopLumpConverter struct{}

func (CatFileChopLumpConverter) ToModels(ethLogs []types.Log) ([]shared.InsertionModel, error) {
	var results []shared.InsertionModel
	for _, ethLog := range ethLogs {
		verifyErr := verifyLog(ethLog)
		if verifyErr != nil {
			return nil, verifyErr
		}
		ilk := ethLog.Topics[2].Hex()
		what := shared.DecodeHexToText(ethLog.Topics[3].Hex())
		dataBytes, parseErr := shared.GetLogNoteArgumentAtIndex(2, ethLog.Data)
		if parseErr != nil {
			return nil, parseErr
		}
		data := shared.ConvertUint256HexToBigInt(hexutil.Encode(dataBytes))

		raw, marshalErr := json.Marshal(ethLog)
		if marshalErr != nil {
			return nil, marshalErr
		}
		result := shared.InsertionModel{
			TableName: "cat_file_chop_lump",
			OrderedColumns: []string{
				"header_id", string(constants2.IlkFK), "what", "data", "tx_idx", "log_idx", "raw_log",
			},
			ColumnValues: shared.ColumnValues{
				"what":    what,
				"data":    data.String(),
				"tx_idx":  ethLog.TxIndex,
				"log_idx": ethLog.Index,
				"raw_log": raw,
			},
			ForeignKeyValues: shared.ForeignKeyValues{
				constants2.IlkFK: ilk,
			},
		}
		results = append(results, result)
	}
	return results, nil
}

func verifyLog(log types.Log) error {
	if len(log.Topics) < 4 {
		return errors.New("log missing topics")
	}
	if len(log.Data) < constants.DataItemLength {
		return errors.New("log missing data")
	}
	return nil
}
