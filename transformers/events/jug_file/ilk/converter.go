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

package ilk

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/pkg/core"
)

type JugFileIlkConverter struct{}

const (
	logDataRequired   = true
	numTopicsRequired = 4
)

func (JugFileIlkConverter) ToModels(_ string, logs []core.HeaderSyncLog) ([]shared.InsertionModel, error) {
	var models []shared.InsertionModel
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

		model := shared.InsertionModel{
			SchemaName: constants.MakerSchema,
			TableName:  constants.JugFileIlkTable,
			OrderedColumns: []string{
				constants.HeaderFK, string(constants.IlkFK), "what", "data", constants.LogFK,
			},
			ColumnValues: shared.ColumnValues{
				"what":             what,
				"data":             data.String(),
				constants.HeaderFK: log.HeaderID,
				constants.LogFK:    log.ID,
			},
			ForeignKeyValues: shared.ForeignKeyValues{
				constants.IlkFK: ilk,
			},
		}
		models = append(models, model)
	}
	return models, nil
}
