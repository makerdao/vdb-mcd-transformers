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
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

type BiteConverter struct{}

func (BiteConverter) toEntities(contractAbi string, logs []core.HeaderSyncLog) ([]BiteEntity, error) {
	var entities []BiteEntity
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

func (converter BiteConverter) ToModels(abi string, logs []core.HeaderSyncLog) ([]shared.InsertionModel, error) {
	entities, entityErr := converter.toEntities(abi, logs)
	if entityErr != nil {
		return nil, fmt.Errorf("BiteConverter couldn't convert logs to entities: %v", entityErr)
	}

	var models []shared.InsertionModel
	for _, biteEntity := range entities {
		ilk := hexutil.Encode(biteEntity.Ilk[:])
		urn := common.BytesToAddress(biteEntity.Urn[:]).Hex()

		model := shared.InsertionModel{
			SchemaName: "maker",
			TableName:  "bite",
			OrderedColumns: []string{
				constants.HeaderFK, constants.LogFK, string(constants.UrnFK), "ink", "art", "tab", "flip", "bid_id",
			},
			ColumnValues: shared.ColumnValues{
				constants.HeaderFK: biteEntity.HeaderID,
				constants.LogFK:    biteEntity.LogID,
				"ink":              shared.BigIntToString(biteEntity.Ink),
				"art":              shared.BigIntToString(biteEntity.Art),
				"tab":              shared.BigIntToString(biteEntity.Tab),
				"flip":             common.BytesToAddress(biteEntity.Flip.Bytes()).Hex(),
				"bid_id":           shared.BigIntToString(biteEntity.Id),
			},
			ForeignKeyValues: shared.ForeignKeyValues{
				constants.IlkFK: ilk,
				constants.UrnFK: urn,
			},
		}
		models = append(models, model)
	}
	return models, nil
}
