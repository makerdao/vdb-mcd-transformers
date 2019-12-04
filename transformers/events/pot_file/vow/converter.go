package vow

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

const (
	logDataRequired                    = false
	numTopicsRequired                  = 4
	What              event.ColumnName = "what"
	Data              event.ColumnName = "data"
)

type Converter struct{}

func (converter Converter) ToModels(_ string, logs []core.HeaderSyncLog, _ *postgres.DB) ([]event.InsertionModel, error) {
	var results []event.InsertionModel
	for _, log := range logs {
		verifyErr := shared.VerifyLog(log.Log, numTopicsRequired, logDataRequired)
		if verifyErr != nil {
			return nil, verifyErr
		}
		what := shared.DecodeHexToText(log.Log.Topics[2].Hex())
		data := common.HexToAddress(log.Log.Topics[3].Hex()).Hex()
		result := event.InsertionModel{
			SchemaName: "maker",
			TableName:  "pot_file_vow",
			OrderedColumns: []event.ColumnName{
				event.HeaderFK,
				event.LogFK,
				What,
				Data,
			},
			ColumnValues: event.ColumnValues{
				event.HeaderFK: log.HeaderID,
				event.LogFK:    log.ID,
				What:           what,
				Data:           data,
			},
		}
		results = append(results, result)
	}
	return results, nil
}
