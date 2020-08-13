package log_bump

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

func (Transformer) toEntities(contractAbi string, logs []core.EventLog) ([]LogBumpEntity, error) {
	var entities []LogBumpEntity
	for _, log := range logs {
		var entity LogBumpEntity
		address := log.Log.Address
		abi, parseErr := eth.ParseAbi(contractAbi)
		if parseErr != nil {
			return nil, parseErr
		}

		contract := bind.NewBoundContract(address, abi, nil, nil, nil)
		unpackErr := contract.UnpackLog(&entity, "LogBump", log.Log)
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
		payGemID, payGemErr := shared.GetOrCreateAddress(entity.PayGem.Hex(), db)
		if payGemErr != nil {
			return nil, shared.ErrCouldNotCreateFK(payGemErr)
		}
		buyGemID, buyGemErr := shared.GetOrCreateAddress(entity.BuyGem.Hex(), db)
		if buyGemErr != nil {
			return nil, shared.ErrCouldNotCreateFK(buyGemErr)
		}
		offerID := big.NewInt(0).SetBytes(entity.Id[:])
		model := event.InsertionModel{
			SchemaName: constants.MakerSchema,
			TableName:  constants.LogBumpTable,
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
				event.HeaderFK:            entity.HeaderID,
				event.LogFK:               entity.LogID,
				event.AddressFK:           addressID,
				constants.OfferId:         shared.BigIntToString(offerID),
				constants.PairColumn:      entity.Pair.Hex(),
				constants.MakerColumn:     makerID,
				constants.PayGemColumn:    payGemID,
				constants.BuyGemColumn:    buyGemID,
				constants.PayAmtColumn:    shared.BigIntToString(entity.PayAmt),
				constants.BuyAmtColumn:    shared.BigIntToString(entity.BuyAmt),
				constants.TimestampColumn: strconv.FormatUint(entity.Timestamp, 10),
			},
		}
		models = append(models, model)
	}
	return models, nil
}
