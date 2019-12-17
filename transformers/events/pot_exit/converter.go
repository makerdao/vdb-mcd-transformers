package pot_exit

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

type Converter struct{}

func (Converter) ToModels(_ string, logs []core.HeaderSyncLog, _ *postgres.DB) ([]event.InsertionModel, error) {
	var models []event.InsertionModel
	for _, log := range logs {
		err := shared.VerifyLog(log.Log, shared.ThreeTopicsRequired, shared.LogDataNotRequired)
		if err != nil {
			return nil, err
		}

		wadInt := shared.ConvertUint256HexToBigInt(log.Log.Topics[2].Hex())

		model := event.InsertionModel{
			SchemaName:     constants.MakerSchema,
			TableName:      constants.PotExitTable,
			OrderedColumns: []event.ColumnName{event.HeaderFK, event.LogFK, constants.WadColumn},
			ColumnValues: event.ColumnValues{
				event.HeaderFK:      log.HeaderID,
				event.LogFK:         log.ID,
				constants.WadColumn: wadInt.String(),
			},
		}
		models = append(models, model)
	}

	return models, nil
}
