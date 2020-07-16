package log_insert

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

func (Transformer) toEntities(contractAbi string, logs []core.EventLog) ([]LogInsertEntity, error) {
	var entities []LogInsertEntity

	for _, log := range logs {
		var entity LogInsertEntity

		abi, parseErr := eth.ParseAbi(contractAbi)
		if parseErr != nil {
			return nil, parseErr
		}
		address := log.Log.Address
		contract := bind.NewBoundContract(address, abi, nil, nil, nil)
		unpackErr := contract.UnpackLog(&entity, "LogInsert", log.Log)
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
		keeperID, keeperErr := shared.GetOrCreateAddress(entity.Keeper.Hex(), db)
		if keeperErr != nil {
			return nil, shared.ErrCouldNotCreateFK(keeperErr)
		}
		addressID, addressErr := shared.GetOrCreateAddress(entity.ContractAddress.Hex(), db)
		if addressErr != nil {
			return nil, shared.ErrCouldNotCreateFK(addressErr)
		}

		model := event.InsertionModel{
			SchemaName: constants.MakerSchema,
			TableName:  constants.LogInsertTable,
			OrderedColumns: []event.ColumnName{
				event.HeaderFK, event.LogFK, event.AddressFK, constants.KeeperColumn,
				constants.OfferId,
			},
			ColumnValues: event.ColumnValues{
				event.AddressFK:        addressID,
				event.HeaderFK:         entity.HeaderID,
				event.LogFK:            entity.LogID,
				constants.KeeperColumn: keeperID,
				constants.OfferId:      shared.BigIntToString(entity.Id),
			},
		}

		models = append(models, model)

	}
	return models, nil
}
