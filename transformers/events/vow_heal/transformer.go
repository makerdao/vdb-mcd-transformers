package vow_heal

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

type Transformer struct{}

func (t Transformer) ToModels(_ string, logs []core.EventLog, db *postgres.DB) ([]event.InsertionModel, error) {
	var models []event.InsertionModel
	for _, log := range logs {
		validationErr := shared.VerifyLog(log.Log, shared.ThreeTopicsRequired, shared.LogDataNotRequired)
		if validationErr != nil {
			return nil, fmt.Errorf("vow heal log failed validation: %w", validationErr)
		}

		msgSenderAddress := common.HexToAddress(log.Log.Topics[1].Hex()).Hex()
		msgSenderAddressID, msgSenderErr := shared.GetOrCreateAddress(msgSenderAddress, db)
		if msgSenderErr != nil {
			msg := "error getting or creating address %s for vow heal msg sender: %w"
			return nil, fmt.Errorf(msg, msgSenderAddress, msgSenderErr)
		}

		radInt := shared.ConvertUint256HexToBigInt(log.Log.Topics[2].Hex())

		model := event.InsertionModel{
			SchemaName: constants.MakerSchema,
			TableName:  constants.VowHealTable,
			OrderedColumns: []event.ColumnName{
				event.HeaderFK, event.LogFK, constants.MsgSenderColumn, constants.RadColumn,
			},
			ColumnValues: event.ColumnValues{
				event.HeaderFK:            log.HeaderID,
				event.LogFK:               log.ID,
				constants.MsgSenderColumn: msgSenderAddressID,
				constants.RadColumn:       radInt.String(),
			},
		}
		models = append(models, model)
	}

	return models, nil
}
