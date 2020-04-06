package batch

import (
	"strconv"

	"github.com/lib/pq"
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

		aLength, aLengthErr := strconv.ParseUint(log.Log.Topics[3].Hex(), 0, 64)
		if aLengthErr != nil {
			return nil, aLengthErr
		}

		aAddresses, aAddressErr := shared.GetLogNoteAddresses(aLength, log.Log.Data)
		if aAddressErr != nil {
			return nil, aAddressErr
		}

		model := event.InsertionModel{
			SchemaName: constants.MakerSchema,
			TableName:  constants.MedianDissBatchTable,
			OrderedColumns: []event.ColumnName{
				event.HeaderFK,
				event.LogFK,
				event.AddressFK,
				constants.MsgSenderColumn,
				constants.ALengthColumn,
				constants.AColumn,
			},
			ColumnValues: event.ColumnValues{
				event.HeaderFK:            log.HeaderID,
				event.LogFK:               log.ID,
				event.AddressFK:           contractAddressID,
				constants.MsgSenderColumn: msgSenderAddressID,
				constants.ALengthColumn:   strconv.FormatUint(aLength, 10),
				constants.AColumn:         pq.Array(aAddresses),
			},
		}
		models = append(models, model)
	}
	return models, nil
}
