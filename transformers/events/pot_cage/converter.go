package pot_cage

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

const (
	logDataRequired   = false
	numTopicsRequired = 2
)

type Converter struct{}

func (converter Converter) ToModels(_ string, logs []core.HeaderSyncLog, _ *postgres.DB) ([]event.InsertionModel, error) {
	var results []event.InsertionModel
	for _, log := range logs {
		verifyErr := shared.VerifyLog(log.Log, numTopicsRequired, logDataRequired)
		if verifyErr != nil {
			return nil, verifyErr
		}
		result := event.InsertionModel{
			SchemaName: "maker",
			TableName:  "pot_cage",
			OrderedColumns: []event.ColumnName{
				event.HeaderFK,
				event.LogFK,
			},
			ColumnValues: event.ColumnValues{
				event.HeaderFK: log.HeaderID,
				event.LogFK:    log.ID,
			},
		}
		results = append(results, result)
	}
	return results, nil
}
