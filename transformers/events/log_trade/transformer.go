package log_trade

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

func (Transformer) toEntities(contractAbi string, logs []core.EventLog) ([]LogTradeEntity, error) {
	var entities []LogTradeEntity
	for _, log := range logs {
		var entity LogTradeEntity
		address := log.Log.Address
		abi, parseErr := eth.ParseAbi(contractAbi)
		if parseErr != nil {
			return nil, parseErr
		}

		contract := bind.NewBoundContract(address, abi, nil, nil, nil)
		unpackErr := contract.UnpackLog(&entity, "LogTrade", log.Log)
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
	for _, entity := range entities {
		addressID, addressErr := shared.GetOrCreateAddress(entity.ContractAddress.Hex(), db)
		if addressErr != nil {
			return nil, shared.ErrCouldNotCreateFK(addressErr)
		}
		payGemID, payGemErr := shared.GetOrCreateAddress(entity.PayGem.Hex(), db)
		if payGemErr != nil {
			return nil, shared.ErrCouldNotCreateFK(payGemErr)
		}
		buyGemID, buyGemErr := shared.GetOrCreateAddress(entity.BuyGem.Hex(), db)
		if buyGemErr != nil {
			return nil, shared.ErrCouldNotCreateFK(buyGemErr)
		}
		model := event.InsertionModel{
			SchemaName: constants.MakerSchema,
			TableName:  constants.LogTradeTable,
			OrderedColumns: []event.ColumnName{
				event.HeaderFK,
				event.LogFK,
				event.AddressFK,
				constants.PayGemColumn,
				constants.BuyGemColumn,
				constants.PayAmtColumn,
				constants.BuyAmtColumn,
			},
			ColumnValues: event.ColumnValues{
				event.HeaderFK:         entity.HeaderID,
				event.LogFK:            entity.LogID,
				event.AddressFK:        addressID,
				constants.PayGemColumn: payGemID,
				constants.BuyGemColumn: buyGemID,
				constants.PayAmtColumn: shared.BigIntToString(entity.PayAmt),
				constants.BuyAmtColumn: shared.BigIntToString(entity.BuyAmt),
			},
		}
		models = append(models, model)
	}
	return models, nil
}
