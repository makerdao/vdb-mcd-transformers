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

package flap_kick

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

func (t Transformer) toEntities(contractAbi string, logs []core.EventLog) ([]FlapKickEntity, error) {
	var entities []FlapKickEntity
	abi, parseErr := eth.ParseAbi(contractAbi)
	if parseErr != nil {
		return nil, parseErr
	}

	for _, log := range logs {
		contract := bind.NewBoundContract(log.Log.Address, abi, nil, nil, nil)
		var entity FlapKickEntity
		unpackErr := contract.UnpackLog(&entity, "Kick", log.Log)
		if unpackErr != nil {
			return nil, unpackErr
		}

		entity.ContractAddress = log.Log.Address
		entity.HeaderID = log.HeaderID
		entity.LogID = log.ID
		entities = append(entities, entity)
	}
	return entities, nil
}

func (t Transformer) ToModels(abi string, logs []core.EventLog, db *postgres.DB) ([]event.InsertionModel, error) {
	entities, entityErr := t.toEntities(abi, logs)
	if entityErr != nil {
		return nil, fmt.Errorf("FlapKick transformer couldn't convert logs to entities: %v", entityErr)
	}

	var models []event.InsertionModel
	for _, flapKickEntity := range entities {
		if flapKickEntity.Id == nil {
			return nil, errors.New("flapKick log ID cannot be nil")
		}
		addressId, addressErr := repository.GetOrCreateAddress(db, flapKickEntity.ContractAddress.Hex())
		if addressErr != nil {
			return nil, shared.ErrCouldNotCreateFK(addressErr)
		}

		model := event.InsertionModel{
			SchemaName: constants.MakerSchema,
			TableName:  constants.FlapKickTable,
			OrderedColumns: []event.ColumnName{
				event.HeaderFK,
				event.LogFK,
				event.AddressFK,
				constants.BidIDColumn,
				constants.LotColumn,
				constants.BidColumn,
			},
			ColumnValues: event.ColumnValues{
				event.HeaderFK:        flapKickEntity.HeaderID,
				event.LogFK:           flapKickEntity.LogID,
				event.AddressFK:       addressId,
				constants.BidIDColumn: flapKickEntity.Id.String(),
				constants.LotColumn:   shared.BigIntToString(flapKickEntity.Lot),
				constants.BidColumn:   shared.BigIntToString(flapKickEntity.Bid),
			},
		}

		models = append(models, model)
	}
	return models, nil
}
