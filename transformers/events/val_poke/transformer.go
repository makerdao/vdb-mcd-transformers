package val_poke

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
		validationErr := shared.VerifyLog(log.Log, shared.ThreeTopicsRequired, shared.LogDataNotRequired)
		if validationErr != nil {
			return nil, validationErr
		}

		contractAddressID, contractAddressErr := shared.GetOrCreateAddress(log.Log.Address.String(), db)
		if contractAddressErr != nil {
			return nil, shared.ErrCouldNotCreateFK(contractAddressErr)
		}

		msgSenderAddress := common.HexToAddress(log.Log.Topics[1].Hex()).Hex()
		msgSenderAddressID, msgSenderAddressErr := shared.GetOrCreateAddress(msgSenderAddress, db)
		if msgSenderAddressErr != nil {
			return nil, shared.ErrCouldNotCreateFK(msgSenderAddressErr)
		}

		wut := shared.DecodeHexToText(log.Log.Topics[2].Hex())

		model := event.InsertionModel{
			SchemaName: constants.MakerSchema,
			TableName:  constants.ValPokeTable,
			OrderedColumns: []event.ColumnName{
				event.HeaderFK,
				event.LogFK,
				event.AddressFK,
				constants.MsgSenderColumn,
				constants.WutColumn,
			},
			ColumnValues: event.ColumnValues{
				event.HeaderFK:            log.HeaderID,
				event.LogFK:               log.ID,
				event.AddressFK:           contractAddressID,
				constants.MsgSenderColumn: msgSenderAddressID,
				constants.WutColumn:       wut,
			},
		}
		models = append(models, model)
	}

	return models, nil
}
