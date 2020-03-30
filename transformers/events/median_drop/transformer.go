package median_drop

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

type Transformer struct{}

func (Transformer) ToModels(_ string, logs []core.EventLog, db *postgres.DB) ([]event.InsertionModel, error) {
	var (
		models []event.InsertionModel
	)
	for _, log := range logs {
		err := shared.VerifyLog(log.Log, shared.FourTopicsRequired, shared.LogDataRequired)
		if err != nil {
			return nil, err
		}

		contractAddressID, contractAddressErr := shared.GetOrCreateAddress(log.Log.Address.String(), db)
		if contractAddressErr != nil {
			return nil, shared.ErrCouldNotCreateFK(contractAddressErr)
		}

		msgSender := log.Log.Topics[1].Hex()
		msgSenderAddressID, msgSenderAddressErr := shared.GetOrCreateAddress(msgSender, db)
		if msgSenderAddressErr != nil {
			return nil, shared.ErrCouldNotCreateFK(msgSenderAddressErr)
		}

		addressIDs, addressErr := shared.GetLogNoteData(2, log.Log.Data, db)
		if addressErr != nil {
			return nil, addressErr
		}

		model := event.InsertionModel{
			SchemaName: constants.MakerSchema,
			TableName:  constants.MedianDropTable,
			OrderedColumns: []event.ColumnName{
				event.HeaderFK, event.LogFK, event.AddressFK, constants.MsgSenderColumn, constants.A0Column, constants.A1Column, constants.A2Column, constants.A3Column,
			},
			ColumnValues: event.ColumnValues{
				event.HeaderFK:            log.HeaderID,
				event.LogFK:               log.ID,
				event.AddressFK:           contractAddressID,
				constants.MsgSenderColumn: msgSenderAddressID,
				constants.A0Column:        addressIDs[0],
				constants.A1Column:        addressIDs[1],
				constants.A2Column:        addressIDs[2],
				constants.A3Column:        addressIDs[3],
			},
		}
		models = append(models, model)
	}
	return models, nil
}
