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
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type VatFileIlkConverter struct{}

const (
	logDataRequired   = false
	numTopicsRequired = 4
)

func (VatFileIlkConverter) ToModels(_ string, logs []core.HeaderSyncLog) ([]shared.InsertionModel, error) {
	//NOTE: the vat contract defines its own custom Note event, rather than relying on DS-Note
	var models []shared.InsertionModel
	for _, log := range logs {
		err := shared.VerifyLog(log.Log, numTopicsRequired, logDataRequired)
		if err != nil {
			return nil, err
		}
		ilk := log.Log.Topics[1].Hex()
		what := shared.DecodeHexToText(log.Log.Topics[2].Hex())
		data := log.Log.Topics[3].Big().String()

		model := shared.InsertionModel{
			SchemaName: "maker",
			TableName:  "vat_file_ilk",
			OrderedColumns: []string{
				constants.HeaderFK, string(constants.IlkFK), "what", "data", constants.LogFK,
			},
			ColumnValues: shared.ColumnValues{
				"what":             what,
				"data":             data,
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
