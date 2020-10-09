package auth

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/repository"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

type Transformer struct {
	TableName event.TableName
}

func (t Transformer) ToModels(_ string, logs []core.EventLog, db *postgres.DB) ([]event.InsertionModel, error) {
	var models []event.InsertionModel
	for _, log := range logs {
		validationErr := shared.VerifyLog(log.Log, shared.ThreeTopicsRequired, shared.LogDataNotRequired)
		if validationErr != nil {
			return nil, validationErr
		}

		contractAddress := log.Log.Address.String()
		contractAddressID, contractAddressErr := repository.GetOrCreateAddress(db, contractAddress)
		if contractAddressErr != nil {
			return nil, shared.ErrCouldNotCreateFK(contractAddressErr)
		}

		msgSenderAddress := common.HexToAddress(log.Log.Topics[1].Hex()).Hex()
		msgSenderAddressID, msgSenderAddressErr := repository.GetOrCreateAddress(db, msgSenderAddress)
		if msgSenderAddressErr != nil {
			return nil, shared.ErrCouldNotCreateFK(msgSenderAddressErr)
		}

		usrAddress := common.HexToAddress(log.Log.Topics[2].Hex()).Hex()
		usrAddressID, usrAddressErr := repository.GetOrCreateAddress(db, usrAddress)
		if usrAddressErr != nil {
			return nil, shared.ErrCouldNotCreateFK(usrAddressErr)
		}

		model := event.InsertionModel{
			SchemaName: constants.MakerSchema,
			TableName:  t.TableName,
			OrderedColumns: []event.ColumnName{
				event.HeaderFK,
				event.LogFK,
				event.AddressFK,
				constants.MsgSenderColumn,
				constants.UsrColumn,
			},
			ColumnValues: event.ColumnValues{
				event.HeaderFK:            log.HeaderID,
				event.LogFK:               log.ID,
				event.AddressFK:           contractAddressID,
				constants.MsgSenderColumn: msgSenderAddressID,
				constants.UsrColumn:       usrAddressID,
			},
		}
		models = append(models, model)
	}

	return models, nil
}
