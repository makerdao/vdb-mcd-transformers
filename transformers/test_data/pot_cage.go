package test_data

import (
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
)

var rawPotCageLog = types.Log{
	Address: common.HexToAddress(PotAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.PotCageSignature()),
		common.HexToHash("0x000000000000000000000000dc127c031e5d6e6e8338f825afd7ebbb01db115d"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e06924500900000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 14764699,
	TxHash:      common.HexToHash("0x6133d54cd52dbb29cd240631b963bf0db3976dbb83290b1b90829a371ab283f4"),
	TxIndex:     4,
	BlockHash:   common.HexToHash("0x774d2ae5736c5936775671e1945d0dde88900417e0c76ce3e50f23d434d6cff6"),
	Index:       5,
	Removed:     false,
}

var PotCageHeaderSyncLog = core.HeaderSyncLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawPotCageLog,
	Transformed: false,
}

var potCageModel = event.InsertionModel{
	SchemaName:     constants.MakerSchema,
	TableName:      constants.PotCageTable,
	OrderedColumns: []event.ColumnName{event.HeaderFK, event.LogFK},
	ColumnValues: event.ColumnValues{
		event.HeaderFK: PotCageHeaderSyncLog.HeaderID,
		event.LogFK:    PotCageHeaderSyncLog.ID,
	},
}

func PotCageModel() event.InsertionModel { return CopyEventModel(potCageModel) }
