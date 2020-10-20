package test_data

import (
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
)

var rawVowHealLog = types.Log{
	Address: common.HexToAddress(VowAddress()),
	Topics: []common.Hash{
		common.HexToHash("0x00000000000000000000000000000000000000000000000000000000f37ac61c"),
		common.HexToHash("0x000000000000000000000000ff490eedbbae058cfc6580c1f965ccf326ff6bd4"),
		common.HexToHash("0x0000000000000000000000092fea48a423025b3f8038dddbdfd5854d25ace5be"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e0f37ac61c0000000000000000000000092fea48a423025b3f8038dddbdfd5854d25ace5be0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 9790175,
	TxHash:      common.HexToHash("0x75d9ccd862106c33e533a3bcf4c04e4ae7fc7ca720080a9f5b2c18cdb81a2e39"),
	TxIndex:     169,
	BlockHash:   common.HexToHash("0x09dc0d99eab212043dc7246df1195049d8eef6b6906dc3533437cc65e83581b6"),
	Index:       152,
	Removed:     false,
}

var VowHealEventLog = core.EventLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawVowHealLog,
	Transformed: false,
}

var vowHealModel = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.VowHealTable,
	OrderedColumns: []event.ColumnName{
		event.HeaderFK,
		event.LogFK,
		constants.MsgSenderColumn,
		constants.RadColumn,
	},
	ColumnValues: event.ColumnValues{
		event.HeaderFK: VowHealEventLog.HeaderID,
		event.LogFK:    VowHealEventLog.ID,
		// msg.sender
		constants.RadColumn: shared.ConvertUint256HexToBigInt(rawVowHealLog.Topics[2].Hex()).String(),
	},
}

func VowHealModel() event.InsertionModel {
	return CopyModel(vowHealModel)
}
