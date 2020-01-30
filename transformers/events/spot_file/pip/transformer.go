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

package pip

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

type Transformer struct{}

func (Transformer) ToModels(_ string, logs []core.EventLog, db *postgres.DB) ([]event.InsertionModel, error) {
	var models []event.InsertionModel
	for _, log := range logs {
		err := shared.VerifyLog(log.Log, shared.FourTopicsRequired, shared.LogDataRequired)
		if err != nil {
			return nil, err
		}

		ilk := log.Log.Topics[2].Hex()
		ilkID, ilkErr := shared.GetOrCreateIlk(ilk, db)
		if ilkErr != nil {
			return nil, shared.ErrCouldNotCreateFK(ilkErr)
		}

		what := shared.DecodeHexToText(log.Log.Topics[3].Hex())

		pipBytes, getErr := shared.GetLogNoteArgumentAtIndex(2, log.Log.Data)
		if getErr != nil {
			return nil, getErr
		}
		pipAddress := common.BytesToAddress(pipBytes)
		addressID, addressErr := shared.GetOrCreateAddress(pipAddress.Hex(), db)
		if addressErr != nil {
			return nil, addressErr
		}

		model := event.InsertionModel{
			SchemaName: constants.MakerSchema,
			TableName:  constants.SpotFilePipTable,
			OrderedColumns: []event.ColumnName{
				event.HeaderFK,
				event.LogFK,
				constants.IlkColumn,
				constants.WhatColumn,
				constants.PipColumn,
			},
			ColumnValues: event.ColumnValues{
				event.HeaderFK:       log.HeaderID,
				event.LogFK:          log.ID,
				constants.IlkColumn:  ilkID,
				constants.WhatColumn: what,
				constants.PipColumn:  addressID,
			},
		}
		models = append(models, model)
	}
	return models, nil
}
