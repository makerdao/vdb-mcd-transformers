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
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

type Transformer struct{}

func (Transformer) ToModels(_ string, logs []core.HeaderSyncLog, db *postgres.DB) ([]event.InsertionModel, error) {
	var models []event.InsertionModel
	for _, log := range logs {
		validateErr := shared.VerifyLog(log.Log, shared.FourTopicsRequired, shared.LogDataRequired)
		if validateErr != nil {
			return nil, validateErr
		}

		addressID, addressErr := shared.GetOrCreateAddress(log.Log.Address.String(), db)
		if addressErr != nil {
			return nil, shared.ErrCouldNotCreateFK(addressErr)
		}

		bidId := log.Log.Topics[2].Big()
		lot := log.Log.Topics[3].Big()
		bidBytes, dataErr := shared.GetLogNoteArgumentAtIndex(2, log.Log.Data)
		if dataErr != nil {
			return nil, dataErr
		}
		bid := shared.ConvertUint256HexToBigInt(hexutil.Encode(bidBytes))

		model := event.InsertionModel{
			SchemaName: constants.MakerSchema,
			TableName:  constants.DentTable,
			OrderedColumns: []event.ColumnName{
				event.HeaderFK, event.AddressFK, event.LogFK, constants.BidIDColumn, constants.LotColumn, constants.BidColumn,
			},
			ColumnValues: event.ColumnValues{
				event.HeaderFK:        log.HeaderID,
				event.LogFK:           log.ID,
				event.AddressFK:       addressID,
				constants.BidIDColumn: bidId.String(),
				constants.LotColumn:   lot.String(),
				constants.BidColumn:   bid.String(),
			},
		}
		models = append(models, model)
	}
	return models, nil
}
