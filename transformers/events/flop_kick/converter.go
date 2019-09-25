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
	"github.com/vulcanize/vulcanizedb/pkg/core"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/vulcanize/vulcanizedb/pkg/geth"

	"github.com/vulcanize/mcd_transformers/transformers/shared"
)

type FlopKickConverter struct{}

func (FlopKickConverter) ToEntities(contractAbi string, logs []core.HeaderSyncLog) ([]interface{}, error) {
	var results []interface{}
	for _, log := range logs {
		var entity Entity
		address := log.Log.Address
		abi, parseErr := geth.ParseAbi(contractAbi)
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

func (FlopKickConverter) ToModels(entities []interface{}) ([]interface{}, error) {
	var results []interface{}
	for _, entity := range entities {
		flopKickEntity, ok := entity.(Entity)
		if !ok {
			return nil, fmt.Errorf("entity of type %T, not %T", entity, Entity{})
		}

		model := Model{
			BidId:           shared.BigIntToString(flopKickEntity.Id),
			Lot:             shared.BigIntToString(flopKickEntity.Lot),
			Bid:             shared.BigIntToString(flopKickEntity.Bid),
			Gal:             flopKickEntity.Gal.String(),
			ContractAddress: flopKickEntity.ContractAddress.Hex(),
			HeaderID:        flopKickEntity.HeaderID,
			LogID:           flopKickEntity.LogID,
		}
		results = append(results, model)
	}

	return results, nil
}
