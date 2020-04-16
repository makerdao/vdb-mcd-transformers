package log_unsorted_offer

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

func (Transformer) toEntities(contractAbi string, logs []core.EventLog) ([]LogUnsortedOfferEntity, error) {
	var entities []LogUnsortedOfferEntity
	for _, log := range logs {
		var entity LogUnsortedOfferEntity
		address := log.Log.Address
		abi, parseErr := eth.ParseAbi(contractAbi)
		if parseErr != nil {
			return nil, parseErr
		}

		contract := bind.NewBoundContract(address, abi, nil, nil, nil)
		unpackErr := contract.UnpackLog(&entity, "LogUnsortedOffer", log.Log)
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
	for _, logUnsortedOfferEntity := range entities {
		addressID, addressErr := shared.GetOrCreateAddress(logUnsortedOfferEntity.ContractAddress.Hex(), db)
		if addressErr != nil {
			return nil, shared.ErrCouldNotCreateFK(addressErr)
		}
		model := event.InsertionModel{
			SchemaName: constants.MakerSchema,
			TableName:  constants.LogUnsortedOfferTable,
			OrderedColumns: []event.ColumnName{
				event.HeaderFK, event.LogFK, event.AddressFK, constants.OfferId,
			},
			ColumnValues: event.ColumnValues{
				event.HeaderFK:    logUnsortedOfferEntity.HeaderID,
				event.LogFK:       logUnsortedOfferEntity.LogID,
				event.AddressFK:   addressID,
				constants.OfferId: shared.BigIntToString(logUnsortedOfferEntity.Id),
			},
		}
		models = append(models, model)
	}
	return models, nil
}
