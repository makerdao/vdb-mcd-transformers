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
	"github.com/vulcanize/vulcanizedb/pkg/core"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/vulcanize/vulcanizedb/pkg/geth"

	"github.com/vulcanize/mcd_transformers/transformers/shared"
)

type FlipKickConverter struct{}

func (FlipKickConverter) ToEntities(contractAbi string, logs []core.HeaderSyncLog) ([]interface{}, error) {
	var entities []interface{}
	for _, log := range logs {
		var entity FlipKickEntity
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
		entity.ContractAddress = address
		entity.HeaderID = log.HeaderID
		entity.LogID = log.ID
		entities = append(entities, entity)
	}

	return entities, nil
}

func (FlipKickConverter) ToModels(entities []interface{}) ([]interface{}, error) {
	var models []interface{}
	for _, entity := range entities {
		flipKickEntity, ok := entity.(FlipKickEntity)
		if !ok {
			return nil, fmt.Errorf("entity of type %T, not %T", entity, FlipKickEntity{})
		}
		if flipKickEntity.Id == nil {
			return nil, errors.New("flip kick bid ID cannot be nil")
		}

		id := flipKickEntity.Id.String()
		lot := shared.BigIntToString(flipKickEntity.Lot)
		bid := shared.BigIntToString(flipKickEntity.Bid)
		tab := shared.BigIntToString(flipKickEntity.Tab)
		usr := flipKickEntity.Usr.String()
		gal := flipKickEntity.Gal.String()
		contractAddress := flipKickEntity.ContractAddress.String()

		model := FlipKickModel{
			BidId:           id,
			Lot:             lot,
			Bid:             bid,
			Tab:             tab,
			Usr:             usr,
			Gal:             gal,
			ContractAddress: contractAddress,
			HeaderID:        flipKickEntity.HeaderID,
			LogID:           flipKickEntity.LogID,
		}
		models = append(models, model)
	}
	return models, nil
}
