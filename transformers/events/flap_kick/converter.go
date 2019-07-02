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
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/geth"

	"github.com/vulcanize/mcd_transformers/transformers/shared"
)

type FlapKickConverter struct{}

func (FlapKickConverter) ToEntities(contractAbi string, ethLogs []types.Log) ([]interface{}, error) {
	var entities []interface{}
	abi, parseErr := geth.ParseAbi(contractAbi)
	if parseErr != nil {
		return nil, parseErr
	}

	for _, ethLog := range ethLogs {
		contract := bind.NewBoundContract(ethLog.Address, abi, nil, nil, nil)

		var entity FlapKickEntity
		unpackErr := contract.UnpackLog(&entity, "Kick", ethLog)
		if unpackErr != nil {
			return nil, unpackErr
		}

		entity.Raw = ethLog
		entity.TransactionIndex = ethLog.TxIndex
		entity.LogIndex = ethLog.Index

		entities = append(entities, entity)
	}
	return entities, nil
}

func (FlapKickConverter) ToModels(entities []interface{}) ([]interface{}, error) {
	var models []interface{}
	for _, entity := range entities {
		flapKickEntity, ok := entity.(FlapKickEntity)
		if !ok {
			return nil, fmt.Errorf("entity of type %T, not %T", entity, FlapKickEntity{})
		}

		if flapKickEntity.Id == nil {
			return nil, errors.New("flapKick log ID cannot be nil")
		}

		rawLog, err := json.Marshal(flapKickEntity.Raw)
		if err != nil {
			return nil, err
		}

		model := FlapKickModel{
			BidId:            flapKickEntity.Id.String(),
			Lot:              shared.BigIntToString(flapKickEntity.Lot),
			Bid:              shared.BigIntToString(flapKickEntity.Bid),
			Gal:              flapKickEntity.Gal.String(),
			ContractAddress:  flapKickEntity.Raw.Address.Hex(),
			LogIndex:         flapKickEntity.LogIndex,
			TransactionIndex: flapKickEntity.TransactionIndex,
			Raw:              rawLog,
		}
		models = append(models, model)
	}
	return models, nil
}
