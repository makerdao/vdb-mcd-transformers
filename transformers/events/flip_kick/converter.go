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
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/core"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/vulcanize/vulcanizedb/pkg/eth"

	"github.com/vulcanize/mcd_transformers/transformers/shared"
)

type FlipKickConverter struct{}

func (FlipKickConverter) toEntities(contractAbi string, logs []core.HeaderSyncLog) ([]FlipKickEntity, error) {
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

func (c FlipKickConverter) ToModels(abi string, logs []core.HeaderSyncLog) ([]shared.InsertionModel, error) {
	entities, entityErr := c.toEntities(abi, logs)
	if entityErr != nil {
		return nil, fmt.Errorf("FlipKickConverter couldn't convert logs to entities: %v", entityErr)
	}
	var models []shared.InsertionModel
	for _, flipKickEntity := range entities {
		if flipKickEntity.Id == nil {
			return nil, errors.New("flip kick bid ID cannot be nil")
		}

		model := shared.InsertionModel{
			SchemaName: "maker",
			TableName:  "flip_kick",
			OrderedColumns: []string{
				constants.HeaderFK, constants.LogFK, "bid_id", "lot", "bid", "tab", "usr", "gal", string(constants.AddressFK),
			},
			ColumnValues: shared.ColumnValues{
				constants.HeaderFK: flipKickEntity.HeaderID,
				constants.LogFK:    flipKickEntity.LogID,
				"bid_id":           flipKickEntity.Id.String(),
				"lot":              shared.BigIntToString(flipKickEntity.Lot),
				"bid":              shared.BigIntToString(flipKickEntity.Bid),
				"tab":              shared.BigIntToString(flipKickEntity.Tab),
				"usr":              flipKickEntity.Usr.String(),
				"gal":              flipKickEntity.Gal.String(),
			},
			ForeignKeyValues: shared.ForeignKeyValues{
				constants.AddressFK: flipKickEntity.ContractAddress.String(),
			},
		}

		models = append(models, model)
	}
	return models, nil
}
