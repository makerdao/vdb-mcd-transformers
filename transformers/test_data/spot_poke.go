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

package test_data

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"math/big"
	"math/rand"
)

var rawEthSpotPokeLog = types.Log{
	Address: common.HexToAddress(SpotAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.SpotPokeSignature()),
	},
	Data:        hexutil.MustDecode("0x434f4c352d410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000013c6d703eec370000000000000000000000000000000000000000000026c69aba83d25409ffca1a"),
	BlockNumber: 11257491,
	TxHash:      common.HexToHash("0x1103dd21f6c43d6f895d41935918119d35d000a109f9353c4575959ba01206bd"),
	TxIndex:     1,
	BlockHash:   common.HexToHash("0x94c99a3b2c77f7cbc401f010c50cd298546cd1084c4d436dde4d22a8b7bbe7e1"),
	Index:       2,
	Removed:     false,
}
var SpotPokeHeaderSyncLog = core.HeaderSyncLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawEthSpotPokeLog,
	Transformed: false,
}

func spot() *big.Int {
	spot := big.Int{}
	spot.SetString("46877063947368421052631578", 10)
	return &spot
}

func SpotPokeModel() shared.InsertionModel { return CopyModel(spotPokeModel) }

var spotPokeModel = shared.InsertionModel{
	SchemaName: "maker",
	TableName:  "spot_poke",
	OrderedColumns: []string{
		constants.HeaderFK, constants.LogFK, string(constants.IlkFK), "value", "spot",
	},
	ColumnValues: shared.ColumnValues{
		constants.HeaderFK: SpotPokeHeaderSyncLog.HeaderID,
		constants.LogFK:    SpotPokeHeaderSyncLog.ID,
		"value":            "89066421500000000.000000",
		"spot":             "46877063947368421052631578",
	},
	ForeignKeyValues: shared.ForeignKeyValues{
		constants.IlkFK: "0x434f4c352d410000000000000000000000000000000000000000000000000000",
	},
}
