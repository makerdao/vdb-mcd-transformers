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
	"github.com/vulcanize/vulcanizedb/libraries/shared/factories/event"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/eth"

	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

type Converter struct {
	db *postgres.DB
}

const (
	Ink  event.ColumnName = "ink"
	Art  event.ColumnName = "art"
	Tab  event.ColumnName = "tab"
	Flip event.ColumnName = "flip"
	Id   event.ColumnName = "bid_id"
)

func (Converter) toEntities(contractAbi string, logs []core.HeaderSyncLog) ([]BiteEntity, error) {
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

		entity.HeaderID = log.HeaderID
		entity.LogID = log.ID
		entities = append(entities, entity)
	}

	return entities, nil
}

func (c Converter) ToModels(abi string, logs []core.HeaderSyncLog) ([]event.InsertionModel, error) {
	entities, entityErr := c.toEntities(abi, logs)
	if entityErr != nil {
		return nil, fmt.Errorf("converter couldn't convert logs to entities: %v", entityErr)
	}

	var models []event.InsertionModel
	for _, biteEntity := range entities {
		hexIlk := hexutil.Encode(biteEntity.Ilk[:])
		urn := common.BytesToAddress(biteEntity.Urn[:]).Hex()

		urnID, urnErr := shared.GetOrCreateUrn(urn, hexIlk, c.db)
		if urnErr != nil {
			return nil, shared.ErrCouldNotCreateFK(urnErr)
		}

		model := event.InsertionModel{
			SchemaName: "maker",
			TableName:  "bite",
			OrderedColumns: []event.ColumnName{
				event.HeaderFK, event.LogFK, constants.UrnColumn, Ink, Art, Tab, Flip, Id,
			},
			ColumnValues: event.ColumnValues{
				event.HeaderFK:      biteEntity.HeaderID,
				event.LogFK:         biteEntity.LogID,
				constants.UrnColumn: urnID,
				Ink:                 shared.BigIntToString(biteEntity.Ink),
				Art:                 shared.BigIntToString(biteEntity.Art),
				Tab:                 shared.BigIntToString(biteEntity.Tab),
				Flip:                common.BytesToAddress(biteEntity.Flip.Bytes()).Hex(),
				Id:                  shared.BigIntToString(biteEntity.Id),
			},
		}
		models = append(models, model)
	}
	return models, nil
}

func (c *Converter) SetDB(db *postgres.DB) {
	c.db = db
}
