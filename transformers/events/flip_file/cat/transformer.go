package cat

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

type Transformer struct{}

func (Transformer) ToModels(_ string, logs []core.EventLog, db *postgres.DB) ([]event.InsertionModel, error) {
	var results []event.InsertionModel
	for _, log := range logs {
		err := shared.VerifyLog(log.Log, shared.FourTopicsRequired, shared.LogDataNotRequired)
		if err != nil {
			return nil, err
		}

		what := shared.DecodeHexToText(log.Log.Topics[2].Hex())

		addressID, addressErr := shared.GetOrCreateAddress(log.Log.Address.Hex(), db)
		if addressErr != nil {
			return nil, shared.ErrCouldNotCreateFK(addressErr)
		}

		msgSenderID, msgSenderErr := shared.GetOrCreateAddress(log.Log.Topics[1].Hex(), db)
		if msgSenderErr != nil {
			return nil, shared.ErrCouldNotCreateFK(msgSenderErr)
		}

		dataColumnID, dataColumnErr := shared.GetOrCreateAddress(log.Log.Topics[3].Hex(), db)
		if dataColumnErr != nil {
			return nil, shared.ErrCouldNotCreateFK(dataColumnErr)
		}

		result := event.InsertionModel{
			SchemaName: constants.MakerSchema,
			TableName:  constants.FlipFileCatTable,
			OrderedColumns: []event.ColumnName{
				event.HeaderFK,
				event.LogFK,
				event.AddressFK,
				constants.MsgSenderColumn,
				constants.WhatColumn,
				constants.DataColumn,
			},
			ColumnValues: event.ColumnValues{
				event.HeaderFK:            log.HeaderID,
				event.LogFK:               log.ID,
				event.AddressFK:           addressID,
				constants.MsgSenderColumn: msgSenderID,
				constants.WhatColumn:      what,
				constants.DataColumn:      dataColumnID,
			},
		}

		results = append(results, result)
	}

	return results, nil
}
