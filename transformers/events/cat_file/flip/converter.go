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

package flip

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/libraries/shared/factories/event"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type Converter struct{
	db *postgres.DB
}

const (
	logDataRequired   = true
	numTopicsRequired = 4
	What event.ColumnName = "what"
	Flip event.ColumnName = "flip"
)

func (converter Converter) ToModels(_ string, logs []core.HeaderSyncLog) ([]event.InsertionModel, error) {
	var results []event.InsertionModel
	for _, log := range logs {
		verifyErr := shared.VerifyLog(log.Log, numTopicsRequired, logDataRequired)
		if verifyErr != nil {
			return nil, verifyErr
		}
		ilk := log.Log.Topics[2].Hex()
		ilkId, ilkErr := shared.GetOrCreateIlk(ilk, converter.db)
		if ilkErr != nil {
			shared.ErrCouldNotCreateFK(ilkErr)
		}

		what := shared.DecodeHexToText(log.Log.Topics[3].Hex())
		flipBytes, parseErr := shared.GetLogNoteArgumentAtIndex(2, log.Log.Data)
		if parseErr != nil {
			return nil, parseErr
		}
		flip := common.BytesToAddress(flipBytes).String()

		result := event.InsertionModel{
			SchemaName: "maker",
			TableName:  "cat_file_flip",
			OrderedColumns: []event.ColumnName{
				event.HeaderFK,
				constants.IlkColumn,
				What,
				Flip,
				event.LogFK,
			},
			ColumnValues:   event.ColumnValues{
				event.HeaderFK: log.HeaderID,
				constants.IlkColumn: ilkId,
				What: what,
				Flip: flip,
				event.LogFK: log.ID,
			},
		}

		results = append(results, result)
	}
	return results, nil
}

func (converter *Converter) SetDB(db *postgres.DB) {
	converter.db = db
}

