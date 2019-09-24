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
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/geth"

	"github.com/vulcanize/mcd_transformers/transformers/shared"
)

type FlapKickConverter struct{}

func (c FlapKickConverter) ToEntities(contractAbi string, ethLogs []types.Log) ([]FlapKickEntity, error) {
	var entities []FlapKickEntity
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

func (c FlapKickConverter) ToModels(abi string, logs []types.Log) ([]shared.InsertionModel, error) {
	entities, entityErr := c.ToEntities(abi, logs)
	if entityErr != nil {
		return nil, fmt.Errorf("FlapKickConverter couldn't convert logs to entities: %v", entityErr)
	}

	var models []shared.InsertionModel
	for _, flapKickEntity := range entities {
		if flapKickEntity.Id == nil {
			return nil, errors.New("flapKick log ID cannot be nil")
		}

		rawLog, err := json.Marshal(flapKickEntity.Raw)
		if err != nil {
			return nil, err
		}

		model := shared.InsertionModel{
			SchemaName: "maker",
			TableName:  "flap_kick",
			OrderedColumns: []string{
				"header_id", "bid_id", "lot", "bid", "address_id", "tx_idx", "log_idx", "raw_log",
			},
			ColumnValues: shared.ColumnValues{
				"bid_id":  flapKickEntity.Id.String(),
				"lot":     shared.BigIntToString(flapKickEntity.Lot),
				"bid":     shared.BigIntToString(flapKickEntity.Bid),
				"log_idx": flapKickEntity.LogIndex,
				"tx_idx":  flapKickEntity.TransactionIndex,
				"raw_log": rawLog,
			},
			ForeignKeyValues: shared.ForeignKeyValues{
				constants.AddressFK: flapKickEntity.Raw.Address.Hex(),
			},
		}

		models = append(models, model)
	}
	return models, nil
}
