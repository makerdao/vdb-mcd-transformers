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

package chop_lump

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

type Converter struct {
	db *postgres.DB
}

func (converter Converter) ToModels(_ string, logs []core.HeaderSyncLog) ([]event.InsertionModel, error) {
	var results []event.InsertionModel
	for _, log := range logs {
		verifyErr := shared.VerifyLog(log.Log, shared.FourTopicsRequired, shared.LogDataRequired)
		if verifyErr != nil {
			return nil, verifyErr
		}
		ilk := log.Log.Topics[2].Hex()
		ilkId, ilkErr := shared.GetOrCreateIlk(ilk, converter.db)
		if ilkErr != nil {
			return nil, shared.ErrCouldNotCreateFK(ilkErr)
		}
		what := shared.DecodeHexToText(log.Log.Topics[3].Hex())
		dataBytes, parseErr := shared.GetLogNoteArgumentAtIndex(2, log.Log.Data)
		if parseErr != nil {
			return nil, parseErr
		}
		data := shared.ConvertUint256HexToBigInt(hexutil.Encode(dataBytes))
		result := event.InsertionModel{
			SchemaName: "maker",
			TableName:  "cat_file_chop_lump",
			OrderedColumns: []event.ColumnName{
				event.HeaderFK,
				constants.IlkColumn,
				constants.WhatColumn,
				constants.DataColumn,
				event.LogFK,
			},
			ColumnValues: event.ColumnValues{
				event.HeaderFK:       log.HeaderID,
				constants.IlkColumn:  ilkId,
				constants.WhatColumn: what,
				constants.DataColumn: data.String(),
				event.LogFK:          log.ID,
			},
		}
		results = append(results, result)
	}
	return results, nil
}

func (converter *Converter) SetDB(db *postgres.DB) {
	converter.db = db
}
