package log_median_price

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/repository"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/eth"
)

type Transformer struct{}

func (Transformer) toEntities(contractAbi string, logs []core.EventLog) ([]LogMedianPriceEntity, error) {
	var entities []LogMedianPriceEntity
	for _, log := range logs {
		var entity LogMedianPriceEntity
		address := log.Log.Address
		abi, parseErr := eth.ParseAbi(contractAbi)
		if parseErr != nil {
			return nil, parseErr
		}

		contract := bind.NewBoundContract(address, abi, nil, nil, nil)
		unpackErr := contract.UnpackLog(&entity, "LogMedianPrice", log.Log)
		if unpackErr != nil {
			return nil, unpackErr
		}

		entity.HeaderID = log.HeaderID
		entity.LogID = log.ID
		entity.ContractAddress = address
		entities = append(entities, entity)
	}

	return entities, nil
}

func (t Transformer) ToModels(abi string, logs []core.EventLog, db *postgres.DB) ([]event.InsertionModel, error) {
	entities, entityErr := t.toEntities(abi, logs)
	if entityErr != nil {
		return nil, fmt.Errorf("LogMedianPrice transformer couldn't convert logs to entities: %v", entityErr)
	}

	var models []event.InsertionModel
	for _, LogMedianPriceEntity := range entities {
		addressID, addressErr := repository.GetOrCreateAddress(db, LogMedianPriceEntity.ContractAddress.Hex())
		if addressErr != nil {
			return nil, shared.ErrCouldNotCreateFK(addressErr)
		}
		LogMedianPriceModel := event.InsertionModel{
			SchemaName: constants.MakerSchema,
			TableName:  constants.LogMedianPriceTable,
			OrderedColumns: []event.ColumnName{
				event.HeaderFK,
				event.LogFK,
				event.AddressFK,
				constants.ValColumn,
				constants.AgeColumn,
			},
			ColumnValues: event.ColumnValues{
				event.HeaderFK:      LogMedianPriceEntity.HeaderID,
				event.LogFK:         LogMedianPriceEntity.LogID,
				event.AddressFK:     addressID,
				constants.ValColumn: shared.BigIntToString(LogMedianPriceEntity.Val),
				constants.AgeColumn: shared.BigIntToString(LogMedianPriceEntity.Age),
			},
		}
		models = append(models, LogMedianPriceModel)
	}

	return models, nil
}
