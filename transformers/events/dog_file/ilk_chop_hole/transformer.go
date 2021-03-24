package ilk_chop_hole

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/hexutil"
	mcdShared "github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/repository"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/eth"
)

type Transformer struct{}

func (t Transformer) toEntities(contractAbi string, logs []core.EventLog) ([]DogFileIlkChopHoleEntity, error) {
	var entities []DogFileIlkChopHoleEntity
	for _, log := range logs {
		var entity DogFileIlkChopHoleEntity
		address := log.Log.Address
		abi, parseErr := eth.ParseAbi(contractAbi)
		if parseErr != nil {
			return nil, parseErr
		}
		contract := bind.NewBoundContract(address, abi, nil, nil, nil)
		unpackErr := contract.UnpackLog(&entity, "FileIlkUint256", log.Log)
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
	for _, dogFileIlkChopHoleEntity := range entities {
		addressId, addressErr := repository.GetOrCreateAddress(db, dogFileIlkChopHoleEntity.ContractAddress.Hex())
		if addressErr != nil {
			return nil, shared.ErrCouldNotCreateFK(addressErr)
		}

		ilkHex := hexutil.Encode(dogFileIlkChopHoleEntity.Ilk[:])

		ilkId, ilkErr := mcdShared.GetOrCreateIlk(ilkHex, db)
		if ilkErr != nil {
			return nil, shared.ErrCouldNotCreateFK(ilkErr)
		}

		what := hexutil.Encode(dogFileIlkChopHoleEntity.What[:])

		model := event.InsertionModel{
			SchemaName: constants.MakerSchema,
			TableName:  constants.DogFileIlkChopHoleTable,
			OrderedColumns: []event.ColumnName{
				event.HeaderFK,
				event.LogFK,
				event.AddressFK,
				constants.IlkColumn,
				constants.WhatColumn,
				constants.DataColumn,
			},
			ColumnValues: event.ColumnValues{
				event.HeaderFK:       dogFileIlkChopHoleEntity.HeaderID,
				event.LogFK:          dogFileIlkChopHoleEntity.LogID,
				event.AddressFK:      addressId,
				constants.IlkColumn:  ilkId,
				constants.WhatColumn: what,
				constants.DataColumn: shared.BigIntToString(dogFileIlkChopHoleEntity.Data),
			},
		}
		models = append(models, model)
	}
	return models, nil
}
