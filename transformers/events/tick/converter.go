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

package tick

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

type Converter struct {
	db *postgres.DB
}

func (c Converter) ToModels(_ string, logs []core.HeaderSyncLog) ([]event.InsertionModel, error) {
	var models []event.InsertionModel
	for _, log := range logs {
		validateErr := shared.VerifyLog(log.Log, shared.ThreeTopicsRequired, shared.LogDataNotRequired)
		if validateErr != nil {
			return nil, validateErr
		}

		addressID, addressErr := shared.GetOrCreateAddress(log.Log.Address.String(), c.db)
		if addressErr != nil {
			return nil, shared.ErrCouldNotCreateFK(addressErr)
		}

		model := event.InsertionModel{
			SchemaName: "maker",
			TableName:  "tick",
			OrderedColumns: []event.ColumnName{
				constants.HeaderFK, constants.LogFK, constants.BidIdColumn, constants.AddressColumn,
			},
			ColumnValues: event.ColumnValues{
				constants.HeaderFK:      log.HeaderID,
				constants.LogFK:         log.ID,
				constants.BidIdColumn:   log.Log.Topics[2].Big().String(),
				constants.AddressColumn: addressID,
			},
		}
		models = append(models, model)
	}
	return models, nil
}

func (c *Converter) SetDB(db *postgres.DB) {
	c.db = db
}
