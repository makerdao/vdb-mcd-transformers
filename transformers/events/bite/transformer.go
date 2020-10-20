// VulcanizeDB
// Copyright © 2018 Vulcanize

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

package bite

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	shared2 "github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/repository"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/eth"
)

type Transformer struct{}

func (Transformer) toEntities(contractAbi string, logs []core.EventLog) ([]BiteEntity, error) {
	var entities []BiteEntity
	for _, log := range logs {
		var entity BiteEntity
		address := log.Log.Address
		abi, parseErr := eth.ParseAbi(contractAbi)
		if parseErr != nil {
			return nil, parseErr
		}

		contract := bind.NewBoundContract(address, abi, nil, nil, nil)
		unpackErr := contract.UnpackLog(&entity, "Bite", log.Log)
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
		return nil, fmt.Errorf("transformer couldn't convert logs to entities: %v", entityErr)
	}

	var models []event.InsertionModel
	for _, biteEntity := range entities {
		hexIlk := hexutil.Encode(biteEntity.Ilk[:])
		urn := common.BytesToAddress(biteEntity.Urn[:]).Hex()

		urnID, urnErr := shared2.GetOrCreateUrn(urn, hexIlk, db)
		if urnErr != nil {
			return nil, shared.ErrCouldNotCreateFK(urnErr)
		}

		addressId, addressErr := repository.GetOrCreateAddress(db, biteEntity.ContractAddress.Hex())
		if addressErr != nil {
			return nil, shared.ErrCouldNotCreateFK(addressErr)
		}

		model := event.InsertionModel{
			SchemaName: constants.MakerSchema,
			TableName:  constants.BiteTable,
			OrderedColumns: []event.ColumnName{
				event.HeaderFK,
				event.LogFK,
				event.AddressFK,
				constants.UrnColumn,
				constants.InkColumn,
				constants.ArtColumn,
				constants.TabColumn,
				constants.FlipColumn,
				constants.BidIDColumn,
			},
			ColumnValues: event.ColumnValues{
				event.HeaderFK:        biteEntity.HeaderID,
				event.LogFK:           biteEntity.LogID,
				event.AddressFK:       addressId,
				constants.UrnColumn:   urnID,
				constants.InkColumn:   shared.BigIntToString(biteEntity.Ink),
				constants.ArtColumn:   shared.BigIntToString(biteEntity.Art),
				constants.TabColumn:   shared.BigIntToString(biteEntity.Tab),
				constants.FlipColumn:  common.BytesToAddress(biteEntity.Flip.Bytes()).Hex(),
				constants.BidIDColumn: shared.BigIntToString(biteEntity.Id),
			},
		}
		models = append(models, model)
	}
	return models, nil
}
