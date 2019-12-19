package cat_rely

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

type Converter struct{}

func (converter Converter) ToModels(_ string, logs []core.HeaderSyncLog, _ *postgres.DB) ([]event.InsertionModel, error) {
	var results []event.InsertionModel
	for _, log := range logs {
		verifyErr := shared.VerifyLog(log.Log, shared.TwoTopicsRequired, shared.LogDataNotRequired)
		if verifyErr != nil {
			return nil, verifyErr
		}

		usr := shared.ConvertUint256HexToBigInt(log.Log.Topics[1].Hex())

		result := event.InsertionModel{
			SchemaName: constants.MakerSchema,
			TableName:  constants.CatRelyTable,
			OrderedColumns: []event.ColumnName{
				event.HeaderFK,
				event.LogFK,
				constants.UsrColumn,
			},
			ColumnValues: event.ColumnValues{
				event.HeaderFK:      log.HeaderID,
				event.LogFK:         log.ID,
				constants.UsrColumn: usr.String(),
			},
		}
		results = append(results, result)
	}
	return results, nil
}
