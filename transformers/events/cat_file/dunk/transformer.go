package dunk

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
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
		err := shared.VerifyLog(log.Log, shared.FourTopicsRequired, shared.LogDataRequired)
		if err != nil {
			return nil, err
		}
		addressID, addressErr := shared.GetOrCreateAddress(log.Log.Address.Hex(), db)
		if addressErr != nil {
			return nil, shared.ErrCouldNotCreateFK(addressErr)
		}
		msgSenderID, msgSenderErr := shared.GetOrCreateAddress(log.Log.Topics[1].Hex(), db)
		if msgSenderErr != nil {
			return nil, shared.ErrCouldNotCreateFK(msgSenderErr)
		}
		ilkID, ilkErr := shared.GetOrCreateIlk(log.Log.Topics[2].Hex(), db)
		if ilkErr != nil {
			return nil, shared.ErrCouldNotCreateFK(ilkErr)
		}
		what := shared.DecodeHexToText(log.Log.Topics[3].Hex())
		dataBytes, parseErr := shared.GetLogNoteArgumentAtIndex(2, log.Log.Data)
		if parseErr != nil {
			return nil, parseErr
		}
		data := shared.ConvertUint256HexToBigInt(hexutil.Encode(dataBytes))

		result := event.InsertionModel{
			SchemaName: constants.MakerSchema,
			TableName:  constants.CatFileDunkTable,
			OrderedColumns: []event.ColumnName{
				event.HeaderFK,
				event.LogFK,
				event.AddressFK,
				constants.MsgSenderColumn,
				constants.IlkColumn,
				constants.WhatColumn,
				constants.DataColumn,
			},
			ColumnValues: event.ColumnValues{
				event.HeaderFK:            log.HeaderID,
				event.LogFK:               log.ID,
				event.AddressFK:           addressID,
				constants.MsgSenderColumn: msgSenderID,
				constants.IlkColumn:       ilkID,
				constants.WhatColumn:      what,
				constants.DataColumn:      data.String(),
			},
		}
		results = append(results, result)
	}
	return results, nil
}
