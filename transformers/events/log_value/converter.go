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

package log_value

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/eth"
)

type Converter struct{}

const Val event.ColumnName = "val"

func (Converter) toEntities(contractAbi string, logs []core.HeaderSyncLog) ([]LogValueEntity, error) {
	var entities []LogValueEntity
	for _, log := range logs {
		var entity LogValueEntity
		address := log.Log.Address
		abi, parseErr := eth.ParseAbi(contractAbi)
		if parseErr != nil {
			return nil, parseErr
		}

		contract := bind.NewBoundContract(address, abi, nil, nil, nil)
		unpackErr := contract.UnpackLog(&entity, "LogValue", log.Log)
		if unpackErr != nil {
			return nil, unpackErr
		}

		entity.HeaderID = log.HeaderID
		entity.LogID = log.ID
		entities = append(entities, entity)
	}

	return entities, nil
}

func (c Converter) ToModels(abi string, logs []core.HeaderSyncLog, _ *postgres.DB) ([]event.InsertionModel, error) {
	entities, entityErr := c.toEntities(abi, logs)
	if entityErr != nil {
		return nil, fmt.Errorf("converter couldn't convert logs to entities: %v", entityErr)
	}

	var models []event.InsertionModel
	for _, logValueEntity := range entities {
		bigIntVal := new(big.Int).SetBytes(logValueEntity.Val[:])

		model := event.InsertionModel{
			SchemaName: "maker",
			TableName:  "log_value",
			OrderedColumns: []event.ColumnName{
				event.HeaderFK, event.LogFK, Val,
			},
			ColumnValues: event.ColumnValues{
				event.HeaderFK: logValueEntity.HeaderID,
				event.LogFK:    logValueEntity.LogID,
				Val:            shared.BigIntToString(bigIntVal),
			},
		}
		models = append(models, model)
	}
	return models, nil
}
