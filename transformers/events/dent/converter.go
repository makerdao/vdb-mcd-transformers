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

package dent

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type DentConverter struct{}

const (
	logDataRequired   = true
	numTopicsRequired = 4
)

func (c DentConverter) ToModels(_ string, logs []core.HeaderSyncLog) (result []shared.InsertionModel, err error) {
	for _, log := range logs {
		validateErr := shared.VerifyLog(log.Log, numTopicsRequired, logDataRequired)
		if validateErr != nil {
			return nil, validateErr
		}

		bidId := log.Log.Topics[2].Big()
		lot := log.Log.Topics[3].Big()
		bidBytes, dataErr := shared.GetLogNoteArgumentAtIndex(2, log.Log.Data)
		if dataErr != nil {
			return nil, dataErr
		}
		bid := shared.ConvertUint256HexToBigInt(hexutil.Encode(bidBytes))

		model := shared.InsertionModel{
			SchemaName: "maker",
			TableName:  "dent",
			OrderedColumns: []string{
				constants.HeaderFK, "bid_id", "lot", "bid", string(constants.AddressFK), constants.LogFK,
			},
			ColumnValues: shared.ColumnValues{
				"bid_id":           bidId.String(),
				"lot":              lot.String(),
				"bid":              bid.String(),
				constants.HeaderFK: log.HeaderID,
				constants.LogFK:    log.ID,
			},
			ForeignKeyValues: shared.ForeignKeyValues{
				constants.AddressFK: log.Log.Address.String(),
			},
		}
		result = append(result, model)
	}
	return result, err
}
