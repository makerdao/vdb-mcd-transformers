package clip_redo

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

func (t Transformer) toEntities(contractAbi string, logs []core.EventLog) ([]ClipRedoEntity, error) {
	var entities []ClipRedoEntity
	for _, log := range logs {
		var entity ClipRedoEntity
		address := log.Log.Address
		abi, parseErr := eth.ParseAbi(contractAbi)
		if parseErr != nil {
			return nil, parseErr
		}
		contract := bind.NewBoundContract(address, abi, nil, nil, nil)
		unpackErr := contract.UnpackLog(&entity, "Redo", log.Log)
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
	for _, ClipRedoEntity := range entities {
		addressId, addressErr := repository.GetOrCreateAddress(db, ClipRedoEntity.ContractAddress.Hex())
		if addressErr != nil {
			return nil, shared.ErrCouldNotCreateFK(addressErr)
		}

		usrId, usrErr := repository.GetOrCreateAddress(db, ClipRedoEntity.Usr.Hex())
		if usrErr != nil {
			return nil, shared.ErrCouldNotCreateFK(usrErr)
		}

		kprId, kprErr := repository.GetOrCreateAddress(db, ClipRedoEntity.Kpr.Hex())
		if kprErr != nil {
			return nil, shared.ErrCouldNotCreateFK(kprErr)
		}

		model := event.InsertionModel{
			SchemaName: constants.MakerSchema,
			TableName:  constants.ClipRedoTable,
			OrderedColumns: []event.ColumnName{
				event.HeaderFK,
				event.LogFK,
				event.AddressFK,
				constants.SaleIDColumn,
				constants.TopColumn,
				constants.TabColumn,
				constants.LotColumn,
				constants.UsrColumn,
				constants.KprColumn,
				constants.CoinColumn,
			},
			ColumnValues: event.ColumnValues{
				event.HeaderFK:         ClipRedoEntity.HeaderID,
				event.LogFK:            ClipRedoEntity.LogID,
				event.AddressFK:        addressId,
				constants.SaleIDColumn: ClipRedoEntity.Id.String(),
				constants.TopColumn:    shared.BigIntToString(ClipRedoEntity.Top),
				constants.TabColumn:    shared.BigIntToString(ClipRedoEntity.Tab),
				constants.LotColumn:    shared.BigIntToString(ClipRedoEntity.Lot),
				constants.UsrColumn:    usrId,
				constants.KprColumn:    kprId,
				constants.CoinColumn:   shared.BigIntToString(ClipRedoEntity.Coin),
			},
		}
		models = append(models, model)
	}

	return models, nil
}
