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

var (
	rawSetMinSellLog = types.Log{
		Address: common.HexToAddress(OasisAddresses()[0]),
		Topics: []common.Hash{
			common.HexToHash(constants.SetMinSellSignature()),
			common.HexToHash("0x000000000000000000000000db33dfd3d61308c33c63209845dad3e6bfb2c674"),
			common.HexToHash("0x0000000000000000000000006b175474e89094c44da98b954eedeac495271d0f"),
			common.HexToHash("0x0000000000000000000000000000000000000000000000001bc16d674ec80000"),
		},

		Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000044bf7c734e0000000000000000000000006b175474e89094c44da98b954eedeac495271d0f0000000000000000000000000000000000000000000000001bc16d674ec80000"),
		BlockNumber: 8944595,
		TxHash:      common.HexToHash("0x307beec267a2ca431f24501b13add423ddfba6aef0b23e4c54673b309f6065cc"),
		TxIndex:     6,
		BlockHash:   fakes.FakeHash,
		Index:       15,
		Removed:     false,
	}

	SetMinSellEventLog = core.EventLog{
		ID:          int64(rand.Int31()),
		HeaderID:    int64(rand.Int31()),
		Log:         rawSetMinSellLog,
		Transformed: false,
	}

	SetMinSellPayGemAddress    = common.HexToAddress("0x6b175474e89094c44da98b954eedeac495271d0f")
	SetMinSellMsgSenderAddress = common.HexToAddress("0xdb33dfd3d61308c33c63209845dad3e6bfb2c674")

	setMinSellModel = event.InsertionModel{
		SchemaName: constants.MakerSchema,
		TableName:  constants.SetMinSellTable,
		OrderedColumns: []event.ColumnName{event.HeaderFK, event.LogFK, event.AddressFK,
			constants.PayGemColumn, constants.MsgSenderColumn, constants.DustColumn},
		ColumnValues: event.ColumnValues{
			event.HeaderFK:       SetMinSellEventLog.HeaderID,
			event.LogFK:          SetMinSellEventLog.ID,
			constants.DustColumn: "2000000000000000000",
		},
	}
)

func SetMinSellModel() event.InsertionModel { return CopyModel(setMinSellModel) }
