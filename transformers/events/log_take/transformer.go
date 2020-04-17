package log_take

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/eth"
)

type Transformer struct{}

func (Transformer) toEntities(contractAbi string, logs []core.EventLog) ([]LogTakeEntity, error) {
	var entities []LogTakeEntity
	for _, log := range logs {
		var entity LogTakeEntity
		address := log.Log.Address
		abi, parseErr := eth.ParseAbi(contractAbi)
		if parseErr != nil {
			return nil, parseErr
		}

		contract := bind.NewBoundContract(address, abi, nil, nil, nil)
		unpackErr := contract.UnpackLog(&entity, "LogTake", log.Log)
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
		makerID, makerErr := shared.GetOrCreateAddress(entity.Maker.Hex(), db)
		if makerErr != nil {
			return nil, shared.ErrCouldNotCreateFK(makerErr)
		}
		takerID, takerErr := shared.GetOrCreateAddress(entity.Taker.Hex(), db)
		if takerErr != nil {
			return nil, shared.ErrCouldNotCreateFK(takerErr)
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
			TableName:  constants.LogTakeTable,
			OrderedColumns: []event.ColumnName{
				event.HeaderFK,
				event.LogFK,
				event.AddressFK,
				constants.OfferId,
				constants.PairColumn,
				constants.MakerColumn,
				constants.PayGemColumn,
				constants.BuyGemColumn,
				constants.TakerColumn,
				constants.TakeAmtColumn,
				constants.GiveAmtColumn,
				constants.TimestampColumn,
			},
			ColumnValues: event.ColumnValues{
				event.HeaderFK:            entity.HeaderID,
				event.LogFK:               entity.LogID,
				event.AddressFK:           addressID,
				constants.OfferId:         shared.BigIntToString(big.NewInt(0).SetBytes(entity.Id.Bytes())),
				constants.PairColumn:      entity.Pair.Hex(),
				constants.MakerColumn:     makerID,
				constants.TakerColumn:     takerID,
				constants.PayGemColumn:    payGemID,
				constants.BuyGemColumn:    buyGemID,
				constants.TakeAmtColumn:   shared.BigIntToString(entity.TakeAmt),
				constants.GiveAmtColumn:   shared.BigIntToString(entity.GiveAmt),
				constants.TimestampColumn: strconv.FormatUint(entity.Timestamp, 10),
			},
		}
		models = append(models, model)
	}
	return models, nil
}
