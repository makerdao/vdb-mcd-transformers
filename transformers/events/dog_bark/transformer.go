package dog_bark

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
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

func (t Transformer) toEntities(contractAbi string, logs []core.EventLog) ([]DogBarkEntity, error) {
	var entities []DogBarkEntity
	for _, log := range logs {
		var entity DogBarkEntity
		address := log.Log.Address
		abi, parseErr := eth.ParseAbi(contractAbi)
		if parseErr != nil {
			return nil, parseErr
		}
		contract := bind.NewBoundContract(address, abi, nil, nil, nil)
		unpackErr := contract.UnpackLog(&entity, "Bark", log.Log)
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
	for _, dogBarkEntity := range entities {
		addressId, addressErr := repository.GetOrCreateAddress(db, dogBarkEntity.ContractAddress.Hex())
		if addressErr != nil {
			return nil, shared.ErrCouldNotCreateFK(addressErr)
		}

		ilkHex := hexutil.Encode(dogBarkEntity.Ilk[:])
		urn := common.BytesToAddress(dogBarkEntity.Urn[:]).Hex()

		urnID, urnErr := mcdShared.GetOrCreateUrn(urn, ilkHex, db)
		if urnErr != nil {
			return nil, shared.ErrCouldNotCreateFK(urnErr)
		}

		ilkId, ilkErr := mcdShared.GetOrCreateIlk(ilkHex, db)
		if ilkErr != nil {
			return nil, shared.ErrCouldNotCreateFK(ilkErr)
		}

		clipAddressID, clipAddressErr := repository.GetOrCreateAddress(db, dogBarkEntity.Clip.Hex())
		if clipAddressErr != nil {
			return nil, shared.ErrCouldNotCreateFK(clipAddressErr)
		}

		model := event.InsertionModel{
			SchemaName: constants.MakerSchema,
			TableName:  constants.DogBarkTable,
			OrderedColumns: []event.ColumnName{
				event.HeaderFK,
				event.LogFK,
				event.AddressFK,
				constants.IlkColumn,
				constants.UrnColumn,
				constants.InkColumn,
				constants.ArtColumn,
				constants.DueColumn,
				constants.ClipColumn,
				constants.SalesIDColumn,
			},
			ColumnValues: event.ColumnValues{
				event.HeaderFK:          dogBarkEntity.HeaderID,
				event.LogFK:             dogBarkEntity.LogID,
				event.AddressFK:         addressId,
				constants.IlkColumn:     ilkId,
				constants.UrnColumn:     urnID,
				constants.InkColumn:     shared.BigIntToString(dogBarkEntity.Ink),
				constants.ArtColumn:     shared.BigIntToString(dogBarkEntity.Art),
				constants.DueColumn:     shared.BigIntToString(dogBarkEntity.Due),
				constants.ClipColumn:    clipAddressID,
				constants.SalesIDColumn: shared.BigIntToString(dogBarkEntity.Id),
			},
		}
		models = append(models, model)
	}

	return models, nil
}
