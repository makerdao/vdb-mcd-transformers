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
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type TickConverter struct{}

const (
	logDataRequired   = false
	numTopicsRequired = 3
)

func (TickConverter) ToModels(_ string, logs []core.HeaderSyncLog) (results []shared.InsertionModel, err error) {
	for _, log := range logs {
		validateErr := shared.VerifyLog(log.Log, numTopicsRequired, logDataRequired)
		if validateErr != nil {
			return nil, validateErr
		}

		model := shared.InsertionModel{
			SchemaName: "maker",
			TableName:  "tick",
			OrderedColumns: []string{
				constants.HeaderFK, "bid_id", string(constants.AddressFK), constants.LogFK,
			},
			ColumnValues: shared.ColumnValues{
				"bid_id":           log.Log.Topics[2].Big().String(),
				constants.HeaderFK: log.HeaderID,
				constants.LogFK:    log.ID,
			},
			ForeignKeyValues: shared.ForeignKeyValues{
				constants.AddressFK: log.Log.Address.String(),
			},
		}
		results = append(results, model)
	}
	return results, err
}
