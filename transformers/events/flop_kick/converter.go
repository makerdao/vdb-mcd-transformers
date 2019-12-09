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

package flop_kick

import (
	"fmt"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"

	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/pkg/core"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/makerdao/vulcanizedb/pkg/eth"

	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
)

type Converter struct{}

func (Converter) toEntities(contractAbi string, logs []core.HeaderSyncLog) ([]FlopKickEntity, error) {
	var results []FlopKickEntity
	for _, log := range logs {
		var entity FlopKickEntity
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
		entity.ContractAddress = log.Log.Address
		entity.HeaderID = log.HeaderID
		entity.LogID = log.ID
		results = append(results, entity)
	}
	return results, nil
}

func (c Converter) ToModels(abi string, logs []core.HeaderSyncLog, db *postgres.DB) ([]event.InsertionModel, error) {
	var results []event.InsertionModel
	entities, entityErr := c.toEntities(abi, logs)
	if entityErr != nil {
		return nil, fmt.Errorf("FlopKick converter couldn't convert logs to entities: %v", entityErr)
	}
	for _, flopKickEntity := range entities {
		addressId, addressErr := shared.GetOrCreateAddress(flopKickEntity.ContractAddress.Hex(), db)
		if addressErr != nil {
			_ = shared.ErrCouldNotCreateFK(addressErr)
		}

		model := event.InsertionModel{
			SchemaName: constants.MakerSchema,
			TableName:  constants.FlopKickTable,
			OrderedColumns: []event.ColumnName{
				event.HeaderFK,
				event.LogFK,
				event.AddressFK,
				constants.BidIdColumn,
				constants.LotColumn,
				constants.BidColumn,
				constants.GalColumn,
			},
			ColumnValues: event.ColumnValues{
				event.HeaderFK:        flopKickEntity.HeaderID,
				event.LogFK:           flopKickEntity.LogID,
				event.AddressFK:       addressId,
				constants.BidIdColumn: shared.BigIntToString(flopKickEntity.Id),
				constants.LotColumn:   shared.BigIntToString(flopKickEntity.Lot),
				constants.BidColumn:   shared.BigIntToString(flopKickEntity.Bid),
				constants.GalColumn:   flopKickEntity.Gal.String(),
			},
		}
		results = append(results, model)
	}

	return results, nil
}
