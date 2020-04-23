package set_min_sell

import (
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
		err := shared.VerifyLog(log.Log, shared.FourTopicsRequired, true)
		if err != nil {
			return nil, err
		}

		addressID, addressErr := shared.GetOrCreateAddress(log.Log.Address.Hex(), db)
		if addressErr != nil {
			return nil, shared.ErrCouldNotCreateFK(addressErr)
		}

		payGemAddress := log.Log.Topics[2].Hex()
		payGemID, payGemErr := shared.GetOrCreateAddress(payGemAddress, db)
		if payGemErr != nil {
			return nil, shared.ErrCouldNotCreateFK(payGemErr)
		}

		dustInt := shared.ConvertUint256HexToBigInt(log.Log.Topics[3].Hex())

		model := event.InsertionModel{
			SchemaName: constants.MakerSchema,
			TableName:  constants.SetMinSellTable,
			OrderedColumns: []event.ColumnName{
				event.HeaderFK, event.LogFK, event.AddressFK, constants.PayGemColumn,
				constants.DustColumn,
			},
			ColumnValues: event.ColumnValues{
				event.AddressFK:        addressID,
				event.HeaderFK:         log.HeaderID,
				event.LogFK:            log.ID,
				constants.PayGemColumn: payGemID,
				constants.DustColumn:   shared.BigIntToString(dustInt),
			},
		}

		models = append(models, model)
	}

	return models, nil
}
