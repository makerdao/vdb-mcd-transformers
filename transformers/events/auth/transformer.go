package auth

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

type Transformer struct {
	LogNoteArgumentOffset int
	TableName             event.TableName
}

func (t Transformer) ToModels(_ string, logs []core.EventLog, db *postgres.DB) ([]event.InsertionModel, error) {
	numTopicsRequired := shared.ThreeTopicsRequired + t.LogNoteArgumentOffset
	usrTopicIndex := 2 + t.LogNoteArgumentOffset
	var models []event.InsertionModel
	for _, log := range logs {
		validationErr := shared.VerifyLog(log.Log, numTopicsRequired, shared.LogDataNotRequired)
		if validationErr != nil {
			return nil, validationErr
		}

		contractAddress := log.Log.Address.String()
		contractAddressID, contractAddressErr := shared.GetOrCreateAddress(contractAddress, db)
		if contractAddressErr != nil {
			return nil, shared.ErrCouldNotCreateFK(contractAddressErr)
		}

		usrAddress := common.HexToAddress(log.Log.Topics[usrTopicIndex].Hex()).Hex()
		usrAddressID, usrAddressErr := shared.GetOrCreateAddress(usrAddress, db)
		if usrAddressErr != nil {
			return nil, shared.ErrCouldNotCreateFK(usrAddressErr)
		}

		model := event.InsertionModel{
			SchemaName:     constants.MakerSchema,
			TableName:      t.TableName,
			OrderedColumns: []event.ColumnName{event.HeaderFK, event.LogFK, event.AddressFK, constants.UsrColumn},
			ColumnValues: event.ColumnValues{
				event.HeaderFK:      log.HeaderID,
				event.LogFK:         log.ID,
				event.AddressFK:     contractAddressID,
				constants.UsrColumn: usrAddressID,
			},
		}
		models = append(models, model)
	}

	return models, nil
}
