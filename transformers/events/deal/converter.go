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
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type DealConverter struct{}

const (
	logDataRequired   = true
	numTopicsRequired = 3
)

func (DealConverter) ToModels(logs []core.HeaderSyncLog) (result []shared.InsertionModel, err error) {
	for _, log := range logs {
		validationErr := shared.VerifyLog(log.Log, numTopicsRequired, logDataRequired)
		if validationErr != nil {
			return nil, validationErr
		}

		bidId := log.Log.Topics[2].Big()

		model := shared.InsertionModel{
			TableName: "deal",
			OrderedColumns: []string{
				constants.HeaderFK, "bid_id", string(constants.AddressFK), constants.LogFK,
			},
			ColumnValues: shared.ColumnValues{
				"bid_id":           bidId.String(),
				constants.HeaderFK: log.HeaderID,
				constants.LogFK:    log.ID,
			},
			ForeignKeyValues: shared.ForeignKeyValues{
				constants.AddressFK: log.Log.Address.String(),
			},
		}
		result = append(result, model)
	}

	return result, nil
}
