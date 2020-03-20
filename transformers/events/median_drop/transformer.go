package median_drop

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

		aBytes, aErr := shared.GetLogNoteArgumentAtIndex(2, log.Log.Data)
		if aErr != nil {
			return nil, aErr
		}
		a1 := common.BytesToAddress(aBytes).String()
		aAddressID, a1AddressErr := shared.GetOrCreateAddress(a1, db)
		if a1AddressErr != nil {
			return nil, shared.ErrCouldNotCreateFK(a1AddressErr)
		}
		a2Bytes, a2Err := shared.GetLogNoteArgumentAtIndex(3, log.Log.Data)
		if a2Err != nil {
			return nil, a2Err
		}
		a2 := common.BytesToAddress(a2Bytes).String()
		a2AddressID, a2AddressErr := shared.GetOrCreateAddress(a2, db)
		if a2AddressErr != nil {
			return nil, shared.ErrCouldNotCreateFK(a2AddressErr)
		}
		a3Bytes, a3Err := shared.GetLogNoteArgumentAtIndex(4, log.Log.Data)
		if a3Err != nil {
			return nil, a3Err
		}
		a3 := common.BytesToAddress(a3Bytes).String()
		a3AddressID, a3AddressErr := shared.GetOrCreateAddress(a3, db)
		if a3AddressErr != nil {
			return nil, shared.ErrCouldNotCreateFK(a3AddressErr)
		}
		a4Bytes, a4Err := shared.GetLogNoteArgumentAtIndex(5, log.Log.Data)
		if a4Err != nil {
			return nil, a4Err
		}
		a4 := common.BytesToAddress(a4Bytes).String()
		a4AddressID, a4AddressErr := shared.GetOrCreateAddress(a4, db)
		if a4AddressErr != nil {
			return nil, shared.ErrCouldNotCreateFK(a4AddressErr)
		}
		model := event.InsertionModel{
			SchemaName: constants.MakerSchema,
			TableName:  constants.MedianDropTable,
			OrderedColumns: []event.ColumnName{
				event.HeaderFK, event.LogFK, event.AddressFK, constants.MsgSenderColumn, constants.AColumn, constants.A2Column, constants.A3Column, constants.A4Column,
			},
			ColumnValues: event.ColumnValues{
				event.HeaderFK:            log.HeaderID,
				event.LogFK:               log.ID,
				event.AddressFK:           contractAddressID,
				constants.MsgSenderColumn: msgSenderAddressID,
				constants.AColumn:         aAddressID,
				constants.A2Column:        a2AddressID,
				constants.A3Column:        a3AddressID,
				constants.A4Column:        a4AddressID,
			},
		}
		models = append(models, model)
	}
	return models, nil
}
