package vow

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

type Transformer struct{}

func (Transformer) ToModels(_ string, logs []core.HeaderSyncLog, _ *postgres.DB) ([]event.InsertionModel, error) {
	var results []event.InsertionModel
	for _, log := range logs {
		verifyErr := shared.VerifyLog(log.Log, shared.FourTopicsRequired, shared.LogDataNotRequired)
		if verifyErr != nil {
			return nil, verifyErr
		}
		what := shared.DecodeHexToText(log.Log.Topics[2].Hex())
		data := common.HexToAddress(log.Log.Topics[3].Hex()).Hex()
		result := event.InsertionModel{
			SchemaName: constants.MakerSchema,
			TableName:  constants.PotFileVowTable,
			OrderedColumns: []event.ColumnName{
				event.HeaderFK,
				event.LogFK,
				constants.WhatColumn,
				constants.DataColumn,
			},
			ColumnValues: event.ColumnValues{
				event.HeaderFK:       log.HeaderID,
				event.LogFK:          log.ID,
				constants.WhatColumn: what,
				constants.DataColumn: data,
			},
		}
		results = append(results, result)
	}
	return results, nil
}
