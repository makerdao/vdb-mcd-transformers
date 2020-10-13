package auction_file

import (
	"fmt"

	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/repository"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

type Transformer struct{}

func (t Transformer) ToModels(_ string, logs []core.EventLog, db *postgres.DB) ([]event.InsertionModel, error) {
	var models []event.InsertionModel
	for _, log := range logs {
		err := shared.VerifyLog(log.Log, shared.FourTopicsRequired, shared.LogDataNotRequired)
		if err != nil {
			return nil, fmt.Errorf("auction file log missing topics: %w", err)
		}

		contractAddressID, contractAddressErr := repository.GetOrCreateAddress(db, log.Log.Address.String())
		if contractAddressErr != nil {
			errorMsgToFmt := "could not get or creat auction contract address id: %w"
			return nil, fmt.Errorf(errorMsgToFmt, shared.ErrCouldNotCreateFK(contractAddressErr))
		}

		msgSender := log.Log.Topics[1].Hex()
		msgSenderAddressID, msgSenderAddressErr := repository.GetOrCreateAddress(db, msgSender)
		if msgSenderAddressErr != nil {
			errorMsgToFmt := "could not get or create msg sender address id: %w"
			return nil, fmt.Errorf(errorMsgToFmt, shared.ErrCouldNotCreateFK(msgSenderAddressErr))
		}

		what := shared.DecodeHexToText(log.Log.Topics[2].Hex())
		data := shared.ConvertUint256HexToBigInt(log.Log.Topics[3].Hex())

		model := event.InsertionModel{
			SchemaName: constants.MakerSchema,
			TableName:  constants.AuctionFileTable,
			OrderedColumns: []event.ColumnName{
				event.HeaderFK,
				event.LogFK,
				event.AddressFK,
				constants.MsgSenderColumn,
				constants.WhatColumn,
				constants.DataColumn,
			},
			ColumnValues: event.ColumnValues{
				event.HeaderFK:            log.HeaderID,
				event.LogFK:               log.ID,
				event.AddressFK:           contractAddressID,
				constants.MsgSenderColumn: msgSenderAddressID,
				constants.WhatColumn:      what,
				constants.DataColumn:      data.String(),
			},
		}
		models = append(models, model)
	}
	return models, nil
}
