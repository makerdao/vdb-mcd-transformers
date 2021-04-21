package vow

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/repository"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/eth"
)

type Transformer struct{}

func (t Transformer) toEntities(contractAbi string, logs []core.EventLog) ([]DogFileVowEntity, error) {
	var entities []DogFileVowEntity
	for _, log := range logs {
		var entity DogFileVowEntity
		address := log.Log.Address
		abi, parseErr := eth.ParseAbi(contractAbi)
		if parseErr != nil {
			return nil, parseErr
		}
		contract := bind.NewBoundContract(address, abi, nil, nil, nil)
		unpackErr := contract.UnpackLog(&entity, "File0", log.Log)
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

func (t Transformer) ToModels(contractAbi string, logs []core.EventLog, db *postgres.DB) ([]event.InsertionModel, error) {
	entities, entityErr := t.toEntities(contractAbi, logs)
	if entityErr != nil {
		return nil, fmt.Errorf("transformer couldn't convert logs to entities: %v", entityErr)
	}

	var models []event.InsertionModel
	for _, dogFileVowEntity := range entities {
		addressId, addressErr := repository.GetOrCreateAddress(db, dogFileVowEntity.ContractAddress.Hex())
		if addressErr != nil {
			return nil, shared.ErrCouldNotCreateFK(addressErr)
		}

		what := hexutil.Encode(dogFileVowEntity.What[:])

		dataAddressId, dataAddressErr := repository.GetOrCreateAddress(db, dogFileVowEntity.Data.Hex())
		if dataAddressErr != nil {
			return nil, shared.ErrCouldNotCreateFK(dataAddressErr)
		}

		model := event.InsertionModel{
			SchemaName: constants.MakerSchema,
			TableName:  constants.DogFileVowTable,
			OrderedColumns: []event.ColumnName{
				event.HeaderFK,
				event.LogFK,
				event.AddressFK,
				constants.WhatColumn,
				constants.DataColumn,
			},
			ColumnValues: event.ColumnValues{
				event.HeaderFK:       dogFileVowEntity.HeaderID,
				event.LogFK:          dogFileVowEntity.LogID,
				event.AddressFK:      addressId,
				constants.WhatColumn: what,
				constants.DataColumn: dataAddressId,
			},
		}
		models = append(models, model)
	}
	return models, nil
}
