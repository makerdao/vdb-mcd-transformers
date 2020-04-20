package log_min_sell

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

func (Transformer) toEntities(contractAbi string, logs []core.EventLog) ([]LogMinSellEntity, error) {
	var entities []LogMinSellEntity

	for _, log := range logs {
		var entity LogMinSellEntity

		abi, parseErr := eth.ParseAbi(contractAbi)
		if parseErr != nil {
			return nil, parseErr
		}
		address := log.Log.Address
		contract := bind.NewBoundContract(address, abi, nil, nil, nil)
		unpackErr := contract.UnpackLog(&entity, "LogMinSell", log.Log)
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
		return nil, fmt.Errorf("transformer couldn't covert logs to entities :%v", entityErr)
	}
	var models []event.InsertionModel
	for _, entity := range entities {
		payGemID, payGemErr := shared.GetOrCreateAddress(entity.PayGem.Hex(), db)
		if payGemErr != nil {
			return nil, shared.ErrCouldNotCreateFK(payGemErr)
		}
		addressID, addressErr := shared.GetOrCreateAddress(entity.ContractAddress.Hex(), db)
		if addressErr != nil {
			return nil, shared.ErrCouldNotCreateFK(addressErr)
		}

		model := event.InsertionModel{
			SchemaName: constants.MakerSchema,
			TableName:  constants.LogMinSellTable,
			OrderedColumns: []event.ColumnName{
				event.HeaderFK, event.LogFK, event.AddressFK, constants.PayGemColumn,
				constants.MinAmountColumn,
			},
			ColumnValues: event.ColumnValues{
				event.AddressFK:           addressID,
				event.HeaderFK:            entity.HeaderID,
				event.LogFK:               entity.LogID,
				constants.PayGemColumn:    payGemID,
				constants.MinAmountColumn: shared.BigIntToString(entity.MinAmount),
			},
		}

		models = append(models, model)

	}
	return models, nil
}
