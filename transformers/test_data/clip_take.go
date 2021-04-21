package test_data

import (
	"math/big"
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
)

var (
	owe   = big.NewInt(1234567)
	max   = big.NewInt(456)
	price = big.NewInt(11234)
)

var rawClipTakeLog = types.Log{
	Address: common.HexToAddress(Clip130Address()),
	Topics: []common.Hash{
		common.HexToHash(constants.ClipTakeSignature()),
		common.HexToHash("0x000000000000000000000000000000000000000000000000000000000012D687"), //id
		common.HexToHash("0x0000000000000000000000007340e006f4135BA6970D43bf43d88DCAD4e7a8CA"), //usr
	},
	Data: hexutil.MustDecode("0x" +
		"00000000000000000000000000000000000000000000000000000000000001c8" + //max
		"0000000000000000000000000000000000000000000000000000000000002BE2" + //price
		"000000000000000000000000000000000000000000000000000000000012D687" + //owe
		"0000000000000000000000000000000000000000000000000000000000000032" + //tab
		"000000000000000000000000000000000000000000000000000000000000000A"), //lot
	BlockNumber: uint64(testBlockNumber),
	TxHash:      common.HexToHash(flipKickTransactionHash),
	TxIndex:     999,
	BlockHash:   fakes.FakeHash,
	Index:       1,
	Removed:     false,
}

var ClipTakeEventLog = core.EventLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawClipTakeLog,
	Transformed: false,
}

var clipTakeModel = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.ClipTakeTable,
	OrderedColumns: []event.ColumnName{
		event.HeaderFK,
		event.LogFK,
		event.AddressFK,
		constants.SaleIDColumn,
		constants.MaxColumn,
		constants.PriceColumn,
		constants.OweColumn,
		constants.TabColumn,
		constants.LotColumn,
		constants.UsrColumn,
	},
	ColumnValues: event.ColumnValues{
		event.HeaderFK: ClipTakeEventLog.HeaderID,
		event.LogFK:    ClipTakeEventLog.ID,
		// event.AddressFK
		constants.SaleIDColumn: ClipTakeEventLog.Log.Topics[1].Big().String(),
		constants.MaxColumn:    max.String(),
		constants.PriceColumn:  price.String(),
		constants.OweColumn:    owe.String(),
		constants.TabColumn:    tab.String(),
		constants.LotColumn:    lot.String(),
		// constants.UsrColumn
	},
}

func ClipTakeModel() event.InsertionModel { return CopyModel(clipTakeModel) }
