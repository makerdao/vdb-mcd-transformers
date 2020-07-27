package par

import (
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
		err := shared.VerifyLog(log.Log, shared.FourTopicsRequired, shared.LogDataNotRequired)
		if err != nil {
			return nil, err
		}

		msgSenderID, msgSenderErr := shared.GetOrCreateAddress(log.Log.Topics[1].Hex(), db)

		if msgSenderErr != nil {
			return nil, shared.ErrCouldNotCreateFK(msgSenderErr)
		}

		what := shared.DecodeHexToText(log.Log.Topics[2].Hex())
		data := shared.ConvertUint256HexToBigInt(log.Log.Topics[3].Hex())

		model := event.InsertionModel{
			SchemaName: constants.MakerSchema,
			TableName:  constants.SpotFileParTable,
			OrderedColumns: []event.ColumnName{
				event.HeaderFK, constants.WhatColumn, constants.DataColumn, event.LogFK, constants.MsgSenderColumn,
			},
			ColumnValues: event.ColumnValues{
				constants.WhatColumn:      what,
				constants.DataColumn:      data.String(),
				event.HeaderFK:            log.HeaderID,
				event.LogFK:               log.ID,
				constants.MsgSenderColumn: msgSenderID,
			},
		}
		models = append(models, model)
	}
	return models, nil
}
