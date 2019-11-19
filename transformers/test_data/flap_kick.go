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

package test_data

import (
	"math/big"
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/fakes"

	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
)

var rawFlapKickLog = types.Log{
	Address:     common.HexToAddress(FlapAddress()),
	Topics:      []common.Hash{common.HexToHash(constants.FlapKickSignature())},
	Data:        hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000003b9aca000000000000000000000000000000000000000000000000000000000001312d00"),
	BlockNumber: 65,
	TxHash:      common.HexToHash("0xee7930b76b6e93974bd3f37824644ae42a89a3887a1131a7bcb3267ab4dc0169"),
	TxIndex:     66,
	BlockHash:   fakes.FakeHash,
	Index:       67,
	Removed:     false,
}

var FlapKickHeaderSyncLog = core.HeaderSyncLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawFlapKickLog,
	Transformed: false,
}

func FlapKickModel() shared.InsertionModel { return CopyModel(flapKickModel) }

var flapKickModel = shared.InsertionModel{
	SchemaName: "maker",
	TableName:  "flap_kick",
	OrderedColumns: []string{
		constants.HeaderFK, constants.LogFK, "bid_id", "lot", "bid", "address_id",
	},
	ColumnValues: shared.ColumnValues{
		constants.HeaderFK: FlapKickHeaderSyncLog.HeaderID,
		constants.LogFK:    FlapKickHeaderSyncLog.ID,
		"bid_id":           big.NewInt(1).String(),
		"lot":              big.NewInt(1000000000).String(),
		"bid":              big.NewInt(20000000).String(),
	},
	ForeignKeyValues: shared.ForeignKeyValues{
		constants.AddressFK: FlapKickHeaderSyncLog.Log.Address.Hex(),
	},
}
