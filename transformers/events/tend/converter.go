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

package tend

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/pkg/core"
)

type TendConverter struct{}

const (
	logDataRequired   = true
	numTopicsRequired = 4
)

func (TendConverter) ToModels(_ string, logs []core.HeaderSyncLog) (results []shared.InsertionModel, err error) {
	for _, log := range logs {
		err := shared.VerifyLog(log.Log, numTopicsRequired, logDataRequired)
		if err != nil {
			return nil, err
		}

		bidId := log.Log.Topics[2].Big()
		lot := log.Log.Topics[3].Big().String()
		rawBid, bidErr := shared.GetLogNoteArgumentAtIndex(2, log.Log.Data)
		if bidErr != nil {
			return nil, bidErr
		}
		bidValue := shared.ConvertUint256HexToBigInt(hexutil.Encode(rawBid)).String()

		model := shared.InsertionModel{
			SchemaName: "maker",
			TableName:  "tend",
			OrderedColumns: []string{
				constants.HeaderFK, "bid_id", "lot", "bid", string(constants.AddressFK), constants.LogFK,
			},
			ColumnValues: shared.ColumnValues{
				"bid_id":           bidId.String(),
				"lot":              lot,
				"bid":              bidValue,
				constants.HeaderFK: log.HeaderID,
				constants.LogFK:    log.ID,
			},
			ForeignKeyValues: shared.ForeignKeyValues{
				constants.AddressFK: log.Log.Address.Hex(),
			},
		}
		results = append(results, model)
	}
	return results, err
}
