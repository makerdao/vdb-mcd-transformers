// VulcanizeDB
// Copyright © 2018 Vulcanize

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

package log_value

import (
	"math/big"

	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

type Transformer struct{}

func (t Transformer) ToModels(_ string, logs []core.EventLog, db *postgres.DB) ([]event.InsertionModel, error) {
	var models []event.InsertionModel
	for _, log := range logs {

		err := shared.VerifyLog(log.Log, shared.OneTopicRequired, shared.LogDataRequired)
		if err != nil {
			return nil, err
		}

		contractAddressID, contractAddressErr := shared.GetOrCreateAddress(log.Log.Address.String(), db)
		if contractAddressErr != nil {
			return nil, err
		}

		bigIntVal := new(big.Int)
		bigIntVal.SetBytes(log.Log.Data)

		model := event.InsertionModel{
			SchemaName: constants.MakerSchema,
			TableName:  constants.LogValueTable,
			OrderedColumns: []event.ColumnName{
				event.HeaderFK, event.LogFK, event.AddressFK, constants.ValColumn,
			},
			ColumnValues: event.ColumnValues{
				event.AddressFK:     contractAddressID,
				event.HeaderFK:      log.HeaderID,
				event.LogFK:         log.ID,
				constants.ValColumn: shared.BigIntToString(bigIntVal),
			},
		}
		models = append(models, model)
	}
	return models, nil
}
