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
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/core"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/vulcanize/vulcanizedb/pkg/eth"

	"github.com/vulcanize/mcd_transformers/transformers/shared"
)

type FlopKickConverter struct{}

func (FlopKickConverter) toEntities(contractAbi string, logs []core.HeaderSyncLog) ([]FlopKickEntity, error) {
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

func (c FlopKickConverter) ToModels(abi string, logs []core.HeaderSyncLog) ([]shared.InsertionModel, error) {
	var results []shared.InsertionModel
	entities, entityErr := c.toEntities(abi, logs)
	if entityErr != nil {
		return nil, fmt.Errorf("FlopKickConverter couldn't convert logs to entities: %v", entityErr)
	}
	for _, flopKickEntity := range entities {
		model := shared.InsertionModel{
			SchemaName: "maker",
			TableName:  "flop_kick",
			OrderedColumns: []string{
				constants.HeaderFK, constants.LogFK, string(constants.AddressFK), "bid_id", "lot", "bid", "gal",
			},
			ColumnValues: shared.ColumnValues{
				constants.HeaderFK: flopKickEntity.HeaderID,
				constants.LogFK:    flopKickEntity.LogID,
				"bid_id":           shared.BigIntToString(flopKickEntity.Id),
				"lot":              shared.BigIntToString(flopKickEntity.Lot),
				"bid":              shared.BigIntToString(flopKickEntity.Bid),
				"gal":              flopKickEntity.Gal.String(),
			},
			ForeignKeyValues: shared.ForeignKeyValues{
				constants.AddressFK: flopKickEntity.ContractAddress.Hex(),
			},
		}
		results = append(results, model)
	}

	return results, nil
}
