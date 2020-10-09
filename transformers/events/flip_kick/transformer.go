// VulcanizeDB
// Copyright © 2019 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package flip_kick

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/repository"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/eth"
)

type Transformer struct{}

func (Transformer) toEntities(contractAbi string, logs []core.EventLog) ([]FlipKickEntity, error) {
	var entities []FlipKickEntity
	for _, log := range logs {
		var entity FlipKickEntity
		address := log.Log.Address
		abi, err := eth.ParseAbi(contractAbi)
		if err != nil {
			return nil, err
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

func (t Transformer) ToModels(abi string, logs []core.EventLog, db *postgres.DB) ([]event.InsertionModel, error) {
	entities, entityErr := t.toEntities(abi, logs)
	if entityErr != nil {
		return nil, fmt.Errorf("FlipKick transformer couldn't convert logs to entities: %v", entityErr)
	}
	var models []event.InsertionModel
	for _, flipKickEntity := range entities {
		if flipKickEntity.Id == nil {
			return nil, errors.New("flip kick bid ID cannot be nil")
		}
		addressId, addressErr := repository.GetOrCreateAddress(db, flipKickEntity.ContractAddress.Hex())
		if addressErr != nil {
			return nil, shared.ErrCouldNotCreateFK(addressErr)
		}

		model := event.InsertionModel{
			SchemaName: constants.MakerSchema,
			TableName:  constants.FlipKickTable,
			OrderedColumns: []event.ColumnName{
				event.HeaderFK,
				event.LogFK,
				event.AddressFK,
				constants.BidIDColumn,
				constants.LotColumn,
				constants.BidColumn,
				constants.TabColumn,
				constants.UsrColumn,
				constants.GalColumn,
			},
			ColumnValues: event.ColumnValues{
				event.HeaderFK:        flipKickEntity.HeaderID,
				event.LogFK:           flipKickEntity.LogID,
				event.AddressFK:       addressId,
				constants.BidIDColumn: flipKickEntity.Id.String(),
				constants.LotColumn:   shared.BigIntToString(flipKickEntity.Lot),
				constants.BidColumn:   shared.BigIntToString(flipKickEntity.Bid),
				constants.TabColumn:   shared.BigIntToString(flipKickEntity.Tab),
				constants.UsrColumn:   flipKickEntity.Usr.String(),
				constants.GalColumn:   flipKickEntity.Gal.String(),
			},
		}

		models = append(models, model)
	}
	return models, nil
}
