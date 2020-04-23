package log_sorted_offer

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/eth"
)

type Transformer struct{}

func (Transformer) toEntities(contractAbi string, logs []core.EventLog) ([]LogSortedOfferEntity, error) {
	var entities []LogSortedOfferEntity
	for _, log := range logs {
		var entity LogSortedOfferEntity
		address := log.Log.Address
		abi, parseErr := eth.ParseAbi(contractAbi)
		if parseErr != nil {
			return nil, parseErr
		}

		contract := bind.NewBoundContract(address, abi, nil, nil, nil)
		unpackErr := contract.UnpackLog(&entity, "LogSortedOffer", log.Log)
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
		return nil, fmt.Errorf("transformer couldn't convert logs to entities: %v", entityErr)
	}
	var models []event.InsertionModel
	for _, logSortedOfferEntity := range entities {
		addressID, addressErr := shared.GetOrCreateAddress(logSortedOfferEntity.ContractAddress.Hex(), db)
		if addressErr != nil {
			return nil, shared.ErrCouldNotCreateFK(addressErr)
		}
		model := event.InsertionModel{
			SchemaName: constants.MakerSchema,
			TableName:  constants.LogSortedOfferTable,
			OrderedColumns: []event.ColumnName{
				event.HeaderFK, event.LogFK, event.AddressFK, constants.OfferId,
			},
			ColumnValues: event.ColumnValues{
				event.HeaderFK:    logSortedOfferEntity.HeaderID,
				event.LogFK:       logSortedOfferEntity.LogID,
				event.AddressFK:   addressID,
				constants.OfferId: shared.BigIntToString(logSortedOfferEntity.Id),
			},
		}
		models = append(models, model)
	}
	return models, nil
}
