//  VulcanizeDB
//  Copyright © 2019 Vulcanize
//
//  This program is free software: you can redistribute it and/or modify
//  it under the terms of the GNU Affero General Public License as published by
//  the Free Software Foundation, either version 3 of the License, or
//  (at your option) any later version.
//
//  This program is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//  GNU Affero General Public License for more details.
//
//  You should have received a copy of the GNU Affero General Public License
//  along with this program.  If not, see <http://www.gnu.org/licenses/>.

package spot_poke

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/eth"
)

type SpotPokeConverter struct{}

func (s SpotPokeConverter) toEntities(contractAbi string, logs []core.HeaderSyncLog) ([]SpotPokeEntity, error) {
	var entities []SpotPokeEntity
	for _, log := range logs {
		var entity SpotPokeEntity
		address := log.Log.Address
		abi, parseErr := eth.ParseAbi(contractAbi)
		if parseErr != nil {
			return nil, parseErr
		}

		contract := bind.NewBoundContract(address, abi, nil, nil, nil)
		unpackErr := contract.UnpackLog(&entity, "Poke", log.Log)
		if unpackErr != nil {
			return nil, unpackErr
		}

		entity.HeaderID = log.HeaderID
		entity.LogID = log.ID
		entities = append(entities, entity)
	}

	return entities, nil
}

func (s SpotPokeConverter) ToModels(abi string, logs []core.HeaderSyncLog) ([]shared.InsertionModel, error) {
	entities, entityErr := s.toEntities(abi, logs)
	if entityErr != nil {
		return nil, fmt.Errorf("NewCDPConverter couldn't convert logs to entities: %v", entityErr)
	}

	var models []shared.InsertionModel
	for _, spotPokeEntity := range entities {
		spotPokeModel := shared.InsertionModel{
			SchemaName: constants.MakerSchema,
			TableName:  constants.SpotPokeTable,
			OrderedColumns: []string{
				constants.HeaderFK, constants.LogFK, string(constants.IlkFK), "value", "spot",
			},
			ColumnValues: shared.ColumnValues{
				constants.HeaderFK: spotPokeEntity.HeaderID,
				constants.LogFK:    spotPokeEntity.LogID,
				"value":            bytesToFloatString(spotPokeEntity.Val[:], 6),
				"spot":             shared.BigIntToString(spotPokeEntity.Spot),
			},
			ForeignKeyValues: shared.ForeignKeyValues{
				constants.IlkFK: hexutil.Encode(spotPokeEntity.Ilk[:]),
			},
		}
		models = append(models, spotPokeModel)
	}

	return models, nil
}

func bytesToFloatString(bytes []byte, precision int) string {
	bigInt := new(big.Int).SetBytes(bytes)
	bigFloat := new(big.Float).SetInt(bigInt)
	return bigFloat.Text('f', precision)
}
