package clip_kick

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

func (t Transformer) toEntities(contractAbi string, logs []core.EventLog) ([]ClipKickEntity, error) {
	var entities []ClipKickEntity
	for _, log := range logs {
		var entity ClipKickEntity
		address := log.Log.Address
		abi, parseErr := eth.ParseAbi(contractAbi)
		if parseErr != nil {
			return nil, parseErr
		}
		contract := bind.NewBoundContract(address, abi, nil, nil, nil)
		unpackErr := contract.UnpackLog(&entity, "Kick", log.Log)
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
	for _, clipKickEntity := range entities {
		addressId, addressErr := repository.GetOrCreateAddress(db, clipKickEntity.ContractAddress.Hex())
		if addressErr != nil {
			return nil, shared.ErrCouldNotCreateFK(addressErr)
		}

		usrId, usrErr := repository.GetOrCreateAddress(db, clipKickEntity.Usr.Hex())
		if usrErr != nil {
			return nil, shared.ErrCouldNotCreateFK(usrErr)
		}

		kprId, kprErr := repository.GetOrCreateAddress(db, clipKickEntity.Kpr.Hex())
		if kprErr != nil {
			return nil, shared.ErrCouldNotCreateFK(kprErr)
		}

		model := event.InsertionModel{
			SchemaName: constants.MakerSchema,
			TableName:  constants.ClipKickTable,
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
				event.HeaderFK:         clipKickEntity.HeaderID,
				event.LogFK:            clipKickEntity.LogID,
				event.AddressFK:        addressId,
				constants.SaleIDColumn: clipKickEntity.Id.String(),
				constants.TopColumn:    shared.BigIntToString(clipKickEntity.Top),
				constants.TabColumn:    shared.BigIntToString(clipKickEntity.Tab),
				constants.LotColumn:    shared.BigIntToString(clipKickEntity.Lot),
				constants.UsrColumn:    usrId,
				constants.KprColumn:    kprId,
				constants.CoinColumn:   shared.BigIntToString(clipKickEntity.Coin),
			},
		}
		models = append(models, model)
	}

	return models, nil
}
