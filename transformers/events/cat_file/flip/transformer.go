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

package flip

import (
	"github.com/ethereum/go-ethereum/common"
	shared2 "github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/repository"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

type Transformer struct{}

func (Transformer) ToModels(_ string, logs []core.EventLog, db *postgres.DB) ([]event.InsertionModel, error) {
	var results []event.InsertionModel
	for _, log := range logs {
		verifyErr := shared.VerifyLog(log.Log, shared.FourTopicsRequired, shared.LogDataRequired)
		if verifyErr != nil {
			return nil, verifyErr
		}
		addressID, addressErr := repository.GetOrCreateAddress(db, log.Log.Address.Hex())
		if addressErr != nil {
			return nil, shared.ErrCouldNotCreateFK(addressErr)
		}
		msgSender := log.Log.Topics[1].Hex()
		msgSenderID, msgSenderErr := repository.GetOrCreateAddress(db, msgSender)
		if msgSenderErr != nil {
			return nil, shared.ErrCouldNotCreateFK(msgSenderErr)
		}
		ilk := log.Log.Topics[2].Hex()
		ilkID, ilkErr := shared2.GetOrCreateIlk(ilk, db)
		if ilkErr != nil {
			return nil, shared.ErrCouldNotCreateFK(ilkErr)
		}
		what := shared.DecodeHexToText(log.Log.Topics[3].Hex())
		flipBytes, parseErr := shared2.GetLogNoteArgumentAtIndex(2, log.Log.Data)
		if parseErr != nil {
			return nil, parseErr
		}
		flip := common.BytesToAddress(flipBytes).String()

		result := event.InsertionModel{
			SchemaName: constants.MakerSchema,
			TableName:  constants.CatFileFlipTable,
			OrderedColumns: []event.ColumnName{
				event.HeaderFK,
				constants.IlkColumn,
				constants.WhatColumn,
				constants.FlipColumn,
				constants.MsgSenderColumn,
				event.LogFK,
				event.AddressFK,
			},
			ColumnValues: event.ColumnValues{
				event.HeaderFK:            log.HeaderID,
				constants.IlkColumn:       ilkID,
				constants.WhatColumn:      what,
				constants.FlipColumn:      flip,
				constants.MsgSenderColumn: msgSenderID,
				event.LogFK:               log.ID,
				event.AddressFK:           addressID,
			},
		}

		results = append(results, result)
	}
	return results, nil
}
