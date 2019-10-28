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

package new_cdp

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/eth"
)

type NewCdpConverter struct{}

func (NewCdpConverter) toEntities(contractAbi string, logs []core.HeaderSyncLog) ([]NewCdpEntity, error) {
	var entities []NewCdpEntity
	for _, log := range logs {
		var entity NewCdpEntity
		address := log.Log.Address
		abi, err := eth.ParseAbi(contractAbi)
		if err != nil {
			return nil, err
		}

		contract := bind.NewBoundContract(address, abi, nil, nil, nil)

		err = contract.UnpackLog(&entity, "NewCdp", log.Log)
		if err != nil {
			return nil, err
		}

		entity.LogID = log.ID
		entity.HeaderID = log.HeaderID

		entities = append(entities, entity)
	}

	return entities, nil
}

func (converter NewCdpConverter) ToModels(abi string, logs []core.HeaderSyncLog) ([]shared.InsertionModel, error) {
	var models []shared.InsertionModel
	entities, entityErr := converter.toEntities(abi, logs)
	if entityErr != nil {
		return nil, fmt.Errorf("NewCDPConverter couldn't convert logs to entities: %v", entityErr)
	}

	for _, newCdpEntity := range entities {
		model := shared.InsertionModel{
			SchemaName: "maker",
			TableName:  "new_cdp",
			OrderedColumns: []string{
				constants.HeaderFK, constants.LogFK, "usr", "own", "cdp",
			},
			ColumnValues: shared.ColumnValues{
				constants.HeaderFK: newCdpEntity.HeaderID,
				constants.LogFK:    newCdpEntity.LogID,
				"usr":              newCdpEntity.Usr.Hex(),
				"own":              newCdpEntity.Own.Hex(),
				"cdp":              shared.BigIntToString(newCdpEntity.Cdp),
			},
			ForeignKeyValues: shared.ForeignKeyValues{},
		}
		models = append(models, model)
	}
	return models, nil
}
