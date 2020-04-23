package log_kill

import (
	"fmt"
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

func (Transformer) toEntities(contractAbi string, logs []core.EventLog) ([]LogKillEntity, error) {
	var entities []LogKillEntity
	for _, log := range logs {
		var entity LogKillEntity
		address := log.Log.Address
		abi, parseErr := eth.ParseAbi(contractAbi)
		if parseErr != nil {
			return nil, parseErr
		}

		contract := bind.NewBoundContract(address, abi, nil, nil, nil)
		unpackErr := contract.UnpackLog(&entity, "LogKill", log.Log)
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
	for _, logKillEntity := range entities {
		contractAddressId, contractAddressErr := shared.GetOrCreateAddress(logKillEntity.ContractAddress.Hex(), db)
		if contractAddressErr != nil {
			return nil, shared.ErrCouldNotCreateFK(contractAddressErr)
		}

		makerAddressId, makerAddressErr := shared.GetOrCreateAddress(logKillEntity.Maker.Hex(), db)
		if makerAddressErr != nil {
			return nil, shared.ErrCouldNotCreateFK(makerAddressErr)
		}

		payGemAddressId, payGemAddressErr := shared.GetOrCreateAddress(logKillEntity.PayGem.Hex(), db)
		if payGemAddressErr != nil {
			return nil, shared.ErrCouldNotCreateFK(payGemAddressErr)
		}

		buyGemAddressId, buyGemAddressErr := shared.GetOrCreateAddress(logKillEntity.BuyGem.Hex(), db)
		if buyGemAddressErr != nil {
			return nil, shared.ErrCouldNotCreateFK(buyGemAddressErr)
		}

		model := event.InsertionModel{
			SchemaName: constants.MakerSchema,
			TableName:  constants.LogKillTable,
			OrderedColumns: []event.ColumnName{
				event.HeaderFK,
				event.LogFK,
				event.AddressFK,
				constants.OfferId,
				constants.PairColumn,
				constants.MakerColumn,
				constants.PayGemColumn,
				constants.BuyGemColumn,
				constants.PayAmtColumn,
				constants.BuyAmtColumn,
				constants.TimestampColumn,
			},
			ColumnValues: event.ColumnValues{
				event.HeaderFK:            logKillEntity.HeaderID,
				event.LogFK:               logKillEntity.LogID,
				event.AddressFK:           contractAddressId,
				constants.OfferId:         shared.BigIntToString(logKillEntity.Id),
				constants.PairColumn:      logKillEntity.Pair.Hex(),
				constants.MakerColumn:     makerAddressId,
				constants.PayGemColumn:    payGemAddressId,
				constants.BuyGemColumn:    buyGemAddressId,
				constants.PayAmtColumn:    shared.BigIntToString(logKillEntity.PayAmt),
				constants.BuyAmtColumn:    shared.BigIntToString(logKillEntity.BuyAmt),
				constants.TimestampColumn: strconv.FormatUint(logKillEntity.Timestamp, 10),
			},
		}
		models = append(models, model)
	}
	return models, nil
}
