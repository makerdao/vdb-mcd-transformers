package auction_address

import (
	"github.com/ethereum/go-ethereum/common"
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
		verifyLogErr := shared.VerifyLog(log.Log, shared.FourTopicsRequired, shared.LogDataNotRequired)
		if verifyLogErr != nil {
			return nil, verifyLogErr
		}

		msgSenderAddress := common.HexToAddress(log.Log.Topics[1].Hex()).Hex()
		msgSenderAddressId, msgSenderAddressErr := shared.GetOrCreateAddress(msgSenderAddress, db)
		if msgSenderAddressErr != nil {
			return nil, msgSenderAddressErr
		}

		what := shared.DecodeHexToText(log.Log.Topics[2].Hex())
		dataAddress := common.HexToAddress(log.Log.Topics[3].Hex()).Hex()
		dataAddressId, dataAddressErr := shared.GetOrCreateAddress(dataAddress, db)
		if dataAddressErr != nil {
			return nil, dataAddressErr
		}

		model := event.InsertionModel{
			SchemaName: constants.MakerSchema,
			TableName:  constants.VowFileAuctionAddressTable,
			OrderedColumns: []event.ColumnName{
				event.HeaderFK, event.LogFK, constants.MsgSenderColumn, constants.WhatColumn, constants.DataColumn,
			},
			ColumnValues: event.ColumnValues{
				event.HeaderFK:            log.HeaderID,
				event.LogFK:               log.ID,
				constants.MsgSenderColumn: msgSenderAddressId,
				constants.WhatColumn:      what,
				constants.DataColumn:      dataAddressId,
			},
		}

		models = append(models, model)
	}
	return models, nil
}
