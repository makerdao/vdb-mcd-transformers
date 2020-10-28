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

package deal

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/repository"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

type Transformer struct{}

func (Transformer) ToModels(_ string, logs []core.EventLog, db *postgres.DB) ([]event.InsertionModel, error) {
	var models []event.InsertionModel
	for _, log := range logs {
		validationErr := shared.VerifyLog(log.Log, shared.ThreeTopicsRequired, shared.LogDataRequired)
		if validationErr != nil {
			return nil, validationErr
		}

		bidId := log.Log.Topics[2].Big()

		addressID, addressErr := repository.GetOrCreateAddress(db, log.Log.Address.String())
		if addressErr != nil {
			return nil, shared.ErrCouldNotCreateFK(addressErr)
		}

		msgSenderID, msgSenderErr := repository.GetOrCreateAddress(db, log.Log.Topics[1].Hex())
		if msgSenderErr != nil {
			return nil, shared.ErrCouldNotCreateFK(msgSenderErr)
		}

		model := event.InsertionModel{
			SchemaName: constants.MakerSchema,
			TableName:  constants.DealTable,
			OrderedColumns: []event.ColumnName{
				event.HeaderFK,
				event.AddressFK,
				event.LogFK,
				constants.MsgSenderColumn,
				constants.BidIDColumn,
			},
			ColumnValues: event.ColumnValues{
				event.HeaderFK:            log.HeaderID,
				event.AddressFK:           addressID,
				event.LogFK:               log.ID,
				constants.MsgSenderColumn: msgSenderID,
				constants.BidIDColumn:     bidId.String(),
			},
		}
		models = append(models, model)
	}

	return models, nil
}
