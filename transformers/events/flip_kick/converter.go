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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/geth"

	"github.com/vulcanize/mcd_transformers/transformers/shared"
)

type FlipKickConverter struct{}

func (FlipKickConverter) ToEntities(contractAbi string, ethLogs []types.Log) ([]FlipKickEntity, error) {
	var entities []FlipKickEntity
	for _, ethLog := range ethLogs {
		entity := &FlipKickEntity{}
		address := ethLog.Address
		abi, err := geth.ParseAbi(contractAbi)
		if err != nil {
			return nil, err
		}

		contract := bind.NewBoundContract(address, abi, nil, nil, nil)

		err = contract.UnpackLog(entity, "Kick", ethLog)
		if err != nil {
			return nil, err
		}
		entity.ContractAddress = address
		entity.Raw = ethLog
		entity.TransactionIndex = ethLog.TxIndex
		entity.LogIndex = ethLog.Index
		entities = append(entities, *entity)
	}

	return entities, nil
}

func (c FlipKickConverter) ToModels(abi string, logs []types.Log) ([]shared.InsertionModel, error) {
	entities, entityErr := c.ToEntities(abi, logs)
	if entityErr != nil {
		return nil, fmt.Errorf("FlipKickConverter couldn't convert logs to entities: %v", entityErr)
	}
	var models []shared.InsertionModel
	for _, flipKickEntity := range entities {
		if flipKickEntity.Id == nil {
			return nil, errors.New("FlipKick log ID cannot be nil.")
		}

		rawLog, err := json.Marshal(flipKickEntity.Raw)
		if err != nil {
			return nil, err
		}

		model := shared.InsertionModel{
			SchemaName:       "maker",
			TableName:        "flip_kick",
			OrderedColumns:   []string{
				"header_id", "bid_id", "lot", "bid", "tab", "usr", "gal", "address_id", "tx_idx", "log_idx", "raw_log",
			},
			ColumnValues:     shared.ColumnValues{
				"bid_id": flipKickEntity.Id.String(),
				"lot": shared.BigIntToString(flipKickEntity.Lot),
				"bid": shared.BigIntToString(flipKickEntity.Bid),
				"tab": shared.BigIntToString(flipKickEntity.Tab),
				"usr": flipKickEntity.Usr.String(),
				"gal": flipKickEntity.Gal.String(),
				"tx_idx": flipKickEntity.TransactionIndex,
				"log_idx": flipKickEntity.LogIndex,
				"raw_log": rawLog,
			},
			ForeignKeyValues: shared.ForeignKeyValues{
				constants.AddressFK: flipKickEntity.ContractAddress.String(),
			},
		}

		models = append(models, model)
	}
	return models, nil
}
