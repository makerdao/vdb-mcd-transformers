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

package pot_join

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

type Converter struct{}

func (Converter) ToModels(_ string, logs []core.HeaderSyncLog, _ *postgres.DB) ([]event.InsertionModel, error) {
	var models []event.InsertionModel
	for _, log := range logs {
		err := shared.VerifyLog(log.Log, shared.ThreeTopicsRequired, shared.LogDataNotRequired)
		if err != nil {
			return nil, err
		}

		wadInt := shared.ConvertUint256HexToBigInt(log.Log.Topics[2].Hex())

		model := event.InsertionModel{
			SchemaName:     constants.MakerSchema,
			TableName:      constants.PotJoinTable,
			OrderedColumns: []event.ColumnName{event.HeaderFK, event.LogFK, constants.WadColumn},
			ColumnValues: event.ColumnValues{
				event.HeaderFK:      log.HeaderID,
				event.LogFK:         log.ID,
				constants.WadColumn: wadInt.String(),
			},
		}
		models = append(models, model)
	}

	return models, nil
}
