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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/rand"
)

var (
	dentData            = "0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e05ff3a382000000000000000000000000000000000000000000000000002386f26fc1000000000000000000000000000000000000000000000000000000470de4df820000000000000000000000000000000000000000000000000000006a94d74f43000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
	dentTransactionHash = "0x5a210319fcd31eea5959fedb4a1b20881c21a21976e23ff19dff3b44cc1c71e8"
	dentBidId           = "10000000000000000"
	dentLot             = "20000000000000000"
	dentBid             = "30000000000000000"
)

var rawDentLog = types.Log{
	Address: common.HexToAddress(EthFlipAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.DentSignature()),
		common.HexToHash("0x00000000000000000000000064d922894153be9eef7b7218dc565d1d0ce2a092"),
		common.HexToHash("0x000000000000000000000000000000000000000000000000002386f26fc10000"),
		common.HexToHash("0x00000000000000000000000000000000000000000000000000470de4df820000"),
	},
	Data:        hexutil.MustDecode(dentData),
	BlockNumber: 15,
	TxHash:      common.HexToHash(dentTransactionHash),
	TxIndex:     5,
	BlockHash:   fakes.FakeHash,
	Index:       2,
	Removed:     false,
}

var DentHeaderSyncLog = core.HeaderSyncLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawDentLog,
	Transformed: false,
}

var DentModel = shared.InsertionModel{
	TableName: "dent",
	OrderedColumns: []string{
		constants.HeaderFK, "bid_id", "lot", "bid", string(constants.AddressFK), constants.LogFK,
	},
	ColumnValues: shared.ColumnValues{
		"bid_id":           dentBidId,
		"lot":              dentLot,
		"bid":              dentBid,
		constants.HeaderFK: DentHeaderSyncLog.HeaderID,
		constants.LogFK:    DentHeaderSyncLog.ID,
	},
	ForeignKeyValues: shared.ForeignKeyValues{
		constants.AddressFK: DentHeaderSyncLog.Log.Address.Hex(),
	},
}
