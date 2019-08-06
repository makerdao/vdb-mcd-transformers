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
	"errors"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type YankConverter struct{}

func (YankConverter) ToModels(logs []core.HeaderSyncLog) (results []shared.InsertionModel, err error) {
	for _, log := range logs {
		validationErr := validateLog(log.Log)
		if validationErr != nil {
			return nil, validationErr
		}

		bidId := log.Log.Topics[2].Big()

		model := shared.InsertionModel{
			TableName: "yank",
			OrderedColumns: []string{
				"header_id", "bid_id", string(constants.AddressFK), "log_id",
			},
			ColumnValues: shared.ColumnValues{
				"bid_id":    bidId.String(),
				"header_id": log.HeaderID,
				"log_id":    log.ID,
			},
			ForeignKeyValues: shared.ForeignKeyValues{
				constants.AddressFK: log.Log.Address.Hex(),
			},
		}
		results = append(results, model)
	}
	return results, err
}

func validateLog(ethLog types.Log) error {
	if len(ethLog.Topics) < 3 {
		return errors.New("yank log does not contain expected topics")
	}
	return nil
}
