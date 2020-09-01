package cat_claw

import (
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
		verifyLogErr := shared.VerifyLog(log.Log, shared.ThreeTopicsRequired, shared.LogDataNotRequired)
		if verifyLogErr != nil {
			return nil, verifyLogErr
		}
		addressID, addressErr := shared.GetOrCreateAddress(log.Log.Address.Hex(), db)
		if addressErr != nil {
			return nil, shared.ErrCouldNotCreateFK(addressErr)
		}
		msgSenderID, msgSenderErr := shared.GetOrCreateAddress(log.Log.Topics[1].Hex(), db)
		if msgSenderErr != nil {
			return nil, shared.ErrCouldNotCreateFK(msgSenderErr)
		}

		rad := shared.ConvertUint256HexToBigInt(log.Log.Topics[2].Hex())
		model := event.InsertionModel{
			SchemaName: constants.MakerSchema,
			TableName:  constants.CatClawTable,
			OrderedColumns: []event.ColumnName{
				event.HeaderFK,
				event.AddressFK,
				constants.MsgSenderColumn,
				event.LogFK,
				constants.RadColumn,
			},
			ColumnValues: event.ColumnValues{
				event.HeaderFK:            log.HeaderID,
				event.AddressFK:           addressID,
				constants.MsgSenderColumn: msgSenderID,
				event.LogFK:               log.ID,
				constants.RadColumn:       rad.String(),
			},
		}

		results = append(results, model)
	}
	return results, nil
}
