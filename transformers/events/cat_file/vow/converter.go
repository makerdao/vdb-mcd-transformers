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

package vow

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

type Converter struct{}

func (Converter) ToModels(_ string, logs []core.HeaderSyncLog, _ *postgres.DB) ([]event.InsertionModel, error) {
	var results []event.InsertionModel
	for _, log := range logs {
		err := shared.VerifyLog(log.Log, shared.FourTopicsRequired, shared.LogDataRequired)
		if err != nil {
			return nil, err
		}

		what := shared.DecodeHexToText(log.Log.Topics[2].Hex())
		data := common.BytesToAddress(log.Log.Topics[3].Bytes()).String()

		result := event.InsertionModel{
			SchemaName: constants.MakerSchema,
			TableName:  constants.CatFileVowTable,
			OrderedColumns: []event.ColumnName{
				constants.HeaderFK, constants.WhatColumn, constants.DataColumn, event.LogFK,
			},
			ColumnValues: event.ColumnValues{
				constants.HeaderFK:   log.HeaderID,
				constants.WhatColumn: what,
				constants.DataColumn: data,
				event.LogFK:          log.ID,
			},
		}
		results = append(results, result)
	}
	return results, nil
}
