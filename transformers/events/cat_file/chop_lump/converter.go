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
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

var (
	chop = "chop"
	lump = "lump"
)

type CatFileChopLumpConverter struct{}

const (
	logDataRequired   = true
	numTopicsRequired = 4
)

func (CatFileChopLumpConverter) ToModels(logs []core.HeaderSyncLog) ([]shared.InsertionModel, error) {
	var results []shared.InsertionModel
	for _, log := range logs {
		verifyErr := shared.VerifyLog(log.Log, numTopicsRequired, logDataRequired)
		if verifyErr != nil {
			return nil, verifyErr
		}
		ilk := log.Log.Topics[2].Hex()
		what := shared.DecodeHexToText(log.Log.Topics[3].Hex())
		dataBytes, parseErr := shared.GetLogNoteArgumentAtIndex(2, log.Log.Data)
		if parseErr != nil {
			return nil, parseErr
		}
		data := shared.ConvertUint256HexToBigInt(hexutil.Encode(dataBytes))

		result := shared.InsertionModel{
			TableName: "cat_file_chop_lump",
			OrderedColumns: []string{
				"header_id", string(constants.IlkFK), "what", "data", "log_id",
			},
			ColumnValues: shared.ColumnValues{
				"what":      what,
				"data":      data.String(),
				"header_id": log.HeaderID,
				"log_id":    log.ID,
			},
			ForeignKeyValues: shared.ForeignKeyValues{
				constants.IlkFK: ilk,
			},
		}
		results = append(results, result)
	}
	return results, nil
}
