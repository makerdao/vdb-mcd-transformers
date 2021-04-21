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

var RawClipRedoLog = types.Log{
	Address: common.HexToAddress(Clip130Address()),
	Topics: []common.Hash{
		common.HexToHash(constants.ClipKickSignature()),
		common.HexToHash("0x000000000000000000000000000000000000000000000000000000000012D687"), //id
		common.HexToHash("0x0000000000000000000000007340e006f4135BA6970D43bf43d88DCAD4e7a8CA"), //usr
		common.HexToHash("0x00000000000000000000000007Fa9eF6609cA7921112231F8f195138ebbA2977"), //kpr
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

var ClipRedoEventLog = core.EventLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         RawClipRedoLog,
	Transformed: false,
}

var clipRedoModel = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.ClipRedoTable,
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
		event.HeaderFK: ClipRedoEventLog.HeaderID,
		event.LogFK:    ClipRedoEventLog.ID,
		// event.AddressFK
		constants.SaleIDColumn: ClipRedoEventLog.Log.Topics[1].Big().String(),
		constants.TopColumn:    top.String(),
		constants.TabColumn:    tab.String(),
		constants.LotColumn:    lot.String(),
		// constants.UsrColumn
		// constants.KprColumn
		constants.CoinColumn: coin.String(),
	},
}

func ClipRedoModel() event.InsertionModel { return CopyModel(clipRedoModel) }
