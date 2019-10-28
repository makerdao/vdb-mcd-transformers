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
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/eth"

	"github.com/vulcanize/mcd_transformers/transformers/shared"
)

type FlapKickConverter struct{}

func (FlapKickConverter) toEntities(contractAbi string, logs []core.HeaderSyncLog) ([]FlapKickEntity, error) {
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

func (c FlapKickConverter) ToModels(abi string, logs []core.HeaderSyncLog) ([]shared.InsertionModel, error) {
	entities, entityErr := c.toEntities(abi, logs)
	if entityErr != nil {
		return nil, fmt.Errorf("FlapKickConverter couldn't convert logs to entities: %v", entityErr)
	}

	var models []shared.InsertionModel
	for _, flapKickEntity := range entities {
		if flapKickEntity.Id == nil {
			return nil, errors.New("flapKick log ID cannot be nil")
		}

		model := shared.InsertionModel{
			SchemaName: "maker",
			TableName:  "flap_kick",
			OrderedColumns: []string{
				constants.HeaderFK, constants.LogFK, "bid_id", "lot", "bid", "address_id",
			},
			ColumnValues: shared.ColumnValues{
				constants.HeaderFK: flapKickEntity.HeaderID,
				constants.LogFK:    flapKickEntity.LogID,
				"bid_id":           flapKickEntity.Id.String(),
				"lot":              shared.BigIntToString(flapKickEntity.Lot),
				"bid":              shared.BigIntToString(flapKickEntity.Bid),
			},
			ForeignKeyValues: shared.ForeignKeyValues{
				constants.AddressFK: flapKickEntity.ContractAddress.Hex(),
			},
		}

		models = append(models, model)
	}
	return models, nil
}
