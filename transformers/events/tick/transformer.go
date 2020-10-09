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
	"github.com/makerdao/vulcanizedb/libraries/shared/repository"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

// Transformer implements the VDB event Transformer interface
type Transformer struct{}

// ToModels transforms log data into general InsertionModels the Repository can persist
func (Transformer) ToModels(_ string, logs []core.EventLog, db *postgres.DB) ([]event.InsertionModel, error) {
	var models []event.InsertionModel
	for _, log := range logs {
		validateErr := shared.VerifyLog(log.Log, shared.ThreeTopicsRequired, shared.LogDataNotRequired)
		if validateErr != nil {
			return nil, validateErr
		}

		addressID, addressErr := repository.GetOrCreateAddress(db, log.Log.Address.String())
		if addressErr != nil {
			return nil, shared.ErrCouldNotCreateFK(addressErr)
		}

		msgSender := shared.GetChecksumAddressString(log.Log.Topics[1].Hex())
		msgSenderID, msgSenderErr := repository.GetOrCreateAddress(db, msgSender)
		if msgSenderErr != nil {
			return nil, shared.ErrCouldNotCreateFK(msgSenderErr)
		}

		model := event.InsertionModel{
			SchemaName: constants.MakerSchema,
			TableName:  constants.TickTable,
			OrderedColumns: []event.ColumnName{
				event.HeaderFK, event.LogFK, constants.BidIDColumn, event.AddressFK, constants.MsgSenderColumn,
			},
			ColumnValues: event.ColumnValues{
				event.HeaderFK:            log.HeaderID,
				event.LogFK:               log.ID,
				constants.BidIDColumn:     log.Log.Topics[2].Big().String(),
				event.AddressFK:           addressID,
				constants.MsgSenderColumn: msgSenderID,
			},
		}
		models = append(models, model)
	}
	return models, nil
}
