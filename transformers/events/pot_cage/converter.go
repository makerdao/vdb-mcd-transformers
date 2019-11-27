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

type Converter struct {
	db *postgres.DB
}

func (converter Converter) ToModels(contractAbi string, logs []core.HeaderSyncLog) ([]event.InsertionModel, error) {
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

func (converter Converter) SetDB(db *postgres.DB) {
	converter.db = db
}
