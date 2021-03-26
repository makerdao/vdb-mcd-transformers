package clip_take

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

func (t Transformer) toEntities(contractAbi string, logs []core.EventLog) ([]ClipTakeEntity, error) {
	var entities []ClipTakeEntity
	for _, log := range logs {
		var entity ClipTakeEntity
		address := log.Log.Address
		abi, parseErr := eth.ParseAbi(contractAbi)
		if parseErr != nil {
			return nil, parseErr
		}
		contract := bind.NewBoundContract(address, abi, nil, nil, nil)
		unpackErr := contract.UnpackLog(&entity, "Take", log.Log)
		if unpackErr != nil {
			return nil, unpackErr
		}

		entity.ContractAddress = address
		entity.HeaderID = log.HeaderID
		entity.LogID = log.ID
		entities = append(entities, entity)
	}
	return entities, nil
}
func (t Transformer) ToModels(contractAbi string, ethLog []core.EventLog, db *postgres.DB) ([]event.InsertionModel, error) {
	entities, entityErr := t.toEntities(contractAbi, ethLog)
	if entityErr != nil {
		return nil, fmt.Errorf("transformer couldn't convert logs to entities: %v", entityErr)
	}
	var models []event.InsertionModel
	for _, clipTakeEntity := range entities {
		addressId, addressErr := repository.GetOrCreateAddress(db, clipTakeEntity.ContractAddress.Hex())
		if addressErr != nil {
			return nil, shared.ErrCouldNotCreateFK(addressErr)
		}

		usrId, usrErr := repository.GetOrCreateAddress(db, clipTakeEntity.Usr.Hex())
		if usrErr != nil {
			return nil, shared.ErrCouldNotCreateFK(usrErr)
		}

		model := event.InsertionModel{
			SchemaName: constants.MakerSchema,
			TableName:  constants.ClipTakeTable,
			OrderedColumns: []event.ColumnName{
				event.HeaderFK,
				event.LogFK,
				event.AddressFK,
				constants.ClipIDColumn,
				constants.MaxColumn,
				constants.PriceColumn,
				constants.OweColumn,
				constants.TabColumn,
				constants.LotColumn,
				constants.UsrColumn,
			},
			ColumnValues: event.ColumnValues{
				event.HeaderFK:         clipTakeEntity.HeaderID,
				event.LogFK:            clipTakeEntity.LogID,
				event.AddressFK:        addressId,
				constants.ClipIDColumn: clipTakeEntity.Id.String(),
				constants.MaxColumn:    shared.BigIntToString(clipTakeEntity.Max),
				constants.PriceColumn:  shared.BigIntToString(clipTakeEntity.Price),
				constants.OweColumn:    shared.BigIntToString(clipTakeEntity.Owe),
				constants.TabColumn:    shared.BigIntToString(clipTakeEntity.Tab),
				constants.LotColumn:    shared.BigIntToString(clipTakeEntity.Lot),
				constants.UsrColumn:    usrId,
			},
		}
		models = append(models, model)
	}

	return models, nil
}
