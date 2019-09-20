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
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/geth"

	"github.com/vulcanize/mcd_transformers/transformers/shared"
)

type BiteConverter struct{}

func (BiteConverter) ToEntities(contractAbi string, logs []core.HeaderSyncLog) ([]interface{}, error) {
	var entities []interface{}
	for _, log := range logs {
		var entity BiteEntity
		address := log.Log.Address
		abi, parseErr := geth.ParseAbi(contractAbi)
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

func (converter BiteConverter) ToModels(entities []interface{}) ([]interface{}, error) {
	var models []interface{}
	for _, entity := range entities {
		biteEntity, ok := entity.(BiteEntity)
		if !ok {
			return nil, fmt.Errorf("entity of type %T, not %T", entity, BiteEntity{})
		}

		ilk := hexutil.Encode(biteEntity.Ilk[:])
		urn := common.BytesToAddress(biteEntity.Urn[:]).Hex()
		ink := biteEntity.Ink
		art := biteEntity.Art
		tab := biteEntity.Tab
		flip := common.BytesToAddress(biteEntity.Flip.Bytes()).Hex()
		id := biteEntity.Id
		logId := biteEntity.LogID

		model := BiteModel{
			Ilk:      ilk,
			Urn:      urn,
			Ink:      shared.BigIntToString(ink),
			Art:      shared.BigIntToString(art),
			Tab:      shared.BigIntToString(tab),
			Flip:     flip,
			Id:       shared.BigIntToString(id),
			HeaderID: biteEntity.HeaderID,
			LogID:    logId,
		}
		models = append(models, model)
	}
	return models, nil
}
