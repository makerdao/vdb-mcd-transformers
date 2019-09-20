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
	"encoding/json"
	"fmt"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/pkg/geth"

	"github.com/vulcanize/mcd_transformers/transformers/shared"
)

type BiteConverter struct{}

func (BiteConverter) ToEntities(contractAbi string, ethLogs []types.Log) ([]BiteEntity, error) {
	var entities []BiteEntity
	for _, ethLog := range ethLogs {
		entity := &BiteEntity{}
		address := ethLog.Address
		abi, err := geth.ParseAbi(contractAbi)
		if err != nil {
			return nil, err
		}

		contract := bind.NewBoundContract(address, abi, nil, nil, nil)

		err = contract.UnpackLog(entity, "Bite", ethLog)
		if err != nil {
			return nil, err
		}

		entity.Raw = ethLog
		entity.LogIndex = ethLog.Index
		entity.TransactionIndex = ethLog.TxIndex

		entities = append(entities, *entity)
	}

	return entities, nil
}

func (converter BiteConverter) ToModels(abi string, logs []types.Log) ([]shared.InsertionModel, error) {
	entities, entityErr := converter.ToEntities(abi, logs)
	if entityErr != nil {
		return nil, fmt.Errorf("BiteConverter couldn't convert logs to entities: %v", entityErr)
	}

	var models []shared.InsertionModel
	for _, biteEntity := range entities {
		ilk := hexutil.Encode(biteEntity.Ilk[:])
		urn := common.BytesToAddress(biteEntity.Urn[:]).Hex()
		rawLog, err := json.Marshal(biteEntity.Raw)
		if err != nil {
			return nil, err
		}

		model := shared.InsertionModel{
			SchemaName: "maker",
			TableName:  "bite",
			OrderedColumns: []string{
				"header_id", string(constants.UrnFK), "ink", "art", "tab", "flip", "bite_identifier", "tx_idx", "log_idx", "raw_log",
			},
			ColumnValues: shared.ColumnValues{
				"ink":             biteEntity.Ink.String(),
				"art":             biteEntity.Art.String(),
				"tab":             biteEntity.Tab.String(),
				"flip":            common.BytesToAddress(biteEntity.Flip.Bytes()).Hex(),
				"bite_identifier": biteEntity.Id.String(),
				"tx_idx":          biteEntity.TransactionIndex,
				"log_idx":         biteEntity.LogIndex,
				"raw_log":         rawLog,
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
