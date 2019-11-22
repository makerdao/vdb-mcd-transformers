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

package yank

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

const (
	logDataRequired   = false
	numTopicsRequired = 3
	Id                = "bid_id"
)

func (c Converter) ToModels(_ string, logs []core.HeaderSyncLog) (results []event.InsertionModel, err error) {
	for _, log := range logs {
		validationErr := shared.VerifyLog(log.Log, numTopicsRequired, logDataRequired)
		if validationErr != nil {
			return nil, validationErr
		}
		addressID, addressErr := shared.GetOrCreateAddress(log.Log.Address.String(), c.db)
		if addressErr != nil {
			return nil, shared.ErrCouldNotCreateFK(addressErr)
		}
		bidId := log.Log.Topics[2].Big()

		model := event.InsertionModel{
			SchemaName: "maker",
			TableName:  "yank",
			OrderedColumns: []event.ColumnName{
				constants.HeaderFK, Id, constants.AddressColumn, constants.LogFK,
			},
			ColumnValues: event.ColumnValues{
				Id:                      bidId.String(),
				constants.HeaderFK:      log.HeaderID,
				constants.LogFK:         log.ID,
				constants.AddressColumn: addressID,
			},
		}
		results = append(results, model)
	}
	return results, err
}

func (c *Converter) SetDB(db *postgres.DB) {
	c.db = db
}
