package test_data

import (
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
)

var rawCatRelyLog = types.Log{
	Address: common.HexToAddress(CatAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.CatRelySignature()),
		common.HexToHash("0x00000000000000000000000039ad5d336a4c08fac74879f796e1ea0af26c1521"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e06924500900000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 14764546,
	TxHash:      fakes.FakeHash,
	TxIndex:     4,
	BlockHash:   fakes.FakeHash,
	Index:       5,
	Removed:     false,
}

var CatRelyHeaderSyncLog = core.HeaderSyncLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawCatRelyLog,
	Transformed: false,
}

var catRelyModel = event.InsertionModel{
	SchemaName:     constants.MakerSchema,
	TableName:      constants.CatRelyTable,
	OrderedColumns: []event.ColumnName{event.HeaderFK, event.LogFK, constants.AddressColumn},
	ColumnValues: event.ColumnValues{
		event.HeaderFK: CatRelyHeaderSyncLog.HeaderID,
		event.LogFK:    CatRelyHeaderSyncLog.ID,
	},
}

func CatRelyModel() event.InsertionModel { return CopyModel(catRelyModel) }
