package par

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

type Converter struct{}

func (c Converter) ToModels(_ string, logs []core.HeaderSyncLog, _ *postgres.DB) ([]event.InsertionModel, error) {
	var models []event.InsertionModel
	for _, log := range logs {
		err := shared.VerifyLog(log.Log, shared.FourTopicsRequired, shared.LogDataNotRequired)
		if err != nil {
			return nil, err
		}

		what := shared.DecodeHexToText(log.Log.Topics[2].Hex())
		data := shared.ConvertUint256HexToBigInt(log.Log.Topics[3].Hex())

		model := event.InsertionModel{
			SchemaName: constants.MakerSchema,
			TableName:  constants.SpotFileParTable,
			OrderedColumns: []event.ColumnName{
				event.HeaderFK, constants.WhatColumn, constants.DataColumn, event.LogFK,
			},
			ColumnValues: event.ColumnValues{
				constants.WhatColumn: what,
				constants.DataColumn: data.String(),
				event.HeaderFK:       log.HeaderID,
				event.LogFK:          log.ID,
			},
		}
		models = append(models, model)
	}
	return models, nil
}
