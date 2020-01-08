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

package vat_fork

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

type Transformer struct{}

func (Transformer) ToModels(_ string, logs []core.EventLog, db *postgres.DB) ([]event.InsertionModel, error) {
	var models []event.InsertionModel
	for _, log := range logs {
		err := shared.VerifyLog(log.Log, shared.FourTopicsRequired, shared.LogDataRequired)
		if err != nil {
			return nil, err
		}

		ilk := log.Log.Topics[1].Hex()
		src := common.BytesToAddress(log.Log.Topics[2].Bytes()).String()
		dst := common.BytesToAddress(log.Log.Topics[3].Bytes()).String()

		ilkID, ilkErr := shared.GetOrCreateIlk(ilk, db)
		if ilkErr != nil {
			return nil, shared.ErrCouldNotCreateFK(ilkErr)
		}

		dinkBytes, dinkErr := shared.GetLogNoteArgumentAtIndex(3, log.Log.Data)
		if dinkErr != nil {
			return nil, dinkErr
		}
		dink := shared.ConvertInt256HexToBigInt(hexutil.Encode(dinkBytes))

		dartBytes, dartErr := shared.GetLogNoteArgumentAtIndex(4, log.Log.Data)
		if dartErr != nil {
			return nil, dartErr
		}
		dart := shared.ConvertInt256HexToBigInt(hexutil.Encode(dartBytes))

		model := event.InsertionModel{
			SchemaName: constants.MakerSchema,
			TableName:  constants.VatForkTable,
			OrderedColumns: []event.ColumnName{
				event.HeaderFK, constants.IlkColumn, constants.SrcColumn, constants.DstColumn, constants.DinkColumn, constants.DartColumn, event.LogFK,
			},
			ColumnValues: event.ColumnValues{
				constants.SrcColumn:  src,
				constants.DstColumn:  dst,
				constants.DinkColumn: dink.String(),
				constants.DartColumn: dart.String(),
				constants.IlkColumn:  ilkID,
				event.HeaderFK:       log.HeaderID,
				event.LogFK:          log.ID,
			},
		}
		models = append(models, model)
	}

	return models, nil
}
