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
	top  = big.NewInt(456)
	coin = big.NewInt(11234)
)

var rawClipKickLog = types.Log{
	Address: common.HexToAddress(ClipLinkAV130Address()),
	Topics: []common.Hash{
		common.HexToHash(constants.ClipKickSignature()),
		common.HexToHash("0x000000000000000000000000000000000000000000000000000000000012D687"),  //id
		common.HexToHash("0x00000000000000000000000007fa9ef6609ca7921112231f8f195138ebba2977"),  //usr
		common.HexToHash("0x000000000000000000000000011b1f2b752a8b674216c07fecb783f64f0cb19e7"), //kpr
	},
	Data: hexutil.MustDecode("0x" +
		"00000000000000000000000000000000000000000000000000000000000001c8" + //top
		"0000000000000000000000000000000000000000000000000000000000000032" + //tab
		"000000000000000000000000000000000000000000000000000000000000000A" + //lot
		"0000000000000000000000000000000000000000000000000000000000002BE2"), //coin
	BlockNumber: uint64(testBlockNumber),
	TxHash:      common.HexToHash(flipKickTransactionHash),
	TxIndex:     999,
	BlockHash:   fakes.FakeHash,
	Index:       1,
	Removed:     false,
}

var ClipKickEventLog = core.EventLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawClipKickLog,
	Transformed: false,
}

var clipKickModel = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.ClipKickTable,
	OrderedColumns: []event.ColumnName{
		event.HeaderFK,
		event.LogFK,
		event.AddressFK,
		constants.SaleIDColumn,
		constants.TopColumn,
		constants.TabColumn,
		constants.LotColumn,
		constants.UsrColumn,
		constants.KprColumn,
		constants.CoinColumn,
	},
	ColumnValues: event.ColumnValues{
		event.HeaderFK: ClipKickEventLog.HeaderID,
		event.LogFK:    ClipKickEventLog.ID,
		// event.AddressFK
		constants.SaleIDColumn: ClipKickEventLog.Log.Topics[1].Big().String(),
		constants.TopColumn:    top.String(),
		constants.TabColumn:    tab.String(),
		constants.LotColumn:    lot.String(),
		// constants.UsrColumn
		// constants.KprColumn
		constants.CoinColumn: coin.String(),
	},
}

func ClipKickModel() event.InsertionModel { return CopyModel(clipKickModel) }
