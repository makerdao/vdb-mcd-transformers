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
	"github.com/vulcanize/vulcanizedb/pkg/core"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
)

type NewCdpConverter struct{}

func (NewCdpConverter) ToEntities(contractAbi string, logs []core.HeaderSyncLog) ([]interface{}, error) {
	var entities []interface{}
	for _, log := range logs {
		entity := &NewCdpEntity{}
		address := log.Log.Address
		abi, err := geth.ParseAbi(contractAbi)
		if err != nil {
			return nil, err
		}

		contract := bind.NewBoundContract(address, abi, nil, nil, nil)

		err = contract.UnpackLog(entity, "NewCdp", log.Log)
		if err != nil {
			return nil, err
		}

		entity.LogID = log.ID
		entity.HeaderID = log.HeaderID

		entities = append(entities, *entity)
	}

	return entities, nil
}

func (converter NewCdpConverter) ToModels(entities []interface{}) ([]interface{}, error) {
	var models []interface{}
	for _, entity := range entities {
		newCdpEntity, ok := entity.(NewCdpEntity)
		if !ok {
			return nil, fmt.Errorf("entity of type %T, not %T", entity, NewCdpEntity{})
		}

		usr := newCdpEntity.Usr.Hex()
		own := newCdpEntity.Own.Hex()
		cdp := newCdpEntity.Cdp

		model := NewCdpModel{
			Usr:      usr,
			Own:      own,
			Cdp:      shared.BigIntToString(cdp),
			LogID:    newCdpEntity.LogID,
			HeaderID: newCdpEntity.HeaderID,
		}
		models = append(models, model)
	}
	return models, nil
}
