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

var (
	rawLogKillLog = types.Log{
		Address: common.HexToAddress(OasisAddresses()[0]),
		Topics: []common.Hash{
			common.HexToHash(constants.LogKillSignature()),
			common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000938be"),
			common.HexToHash("0xd257ccbe93e550a27236e8cc4971336f6cd2d53037ad567f10fbcc28df6a1eb1"),
			common.HexToHash("0x000000000000000000000000fd10abc3506d8a1112935bf7d13a2e39ca7cadbb"),
		},
		Data:        hexutil.MustDecode("0x0000000000000000000000006b175474e89094c44da98b954eedeac495271d0f000000000000000000000000a0b86991c6218b36c1d19d4a2e9eb0ce3606eb480000000000000000000000000000000000000000000000001bc16d674ec800000000000000000000000000000000000000000000000000000000105ef37c9b80000000000000000000000000000000000000000000000000000000005e61635a"),
		BlockNumber: 9613377,
		TxHash:      common.HexToHash("0xbb6226a053c2bd699d44e92dc13a038cafbafc1db0ad73f136e8f6310e7824d3"),
		TxIndex:     0,
		BlockHash:   common.HexToHash("0x05a331993245c0d5593ba0ae0724a1b4161be11834c5f9957c9dc6bff270cb58"),
		Index:       0,
		Removed:     false,
	}

	LogKillEventLog = core.EventLog{
		ID:          rand.Int63(),
		HeaderID:    rand.Int63(),
		Log:         rawLogKillLog,
		Transformed: false,
	}

	PayGemAddress = common.HexToAddress("0x6b175474e89094c44da98b954eedeac495271d0f")
	BuyGemAddress = common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")

	logKillModel = event.InsertionModel{
		SchemaName: constants.MakerSchema,
		TableName:  constants.LogKillTable,
		OrderedColumns: []event.ColumnName{
			event.HeaderFK,
			event.LogFK,
			event.AddressFK,
			constants.OfferId,
			constants.PairColumn,
			constants.MakerColumn,
			constants.PayGemColumn,
			constants.BuyGemColumn,
			constants.PayAmtColumn,
			constants.BuyAmtColumn,
			constants.TimestampColumn,
		},
		ColumnValues: event.ColumnValues{
			event.HeaderFK: LogKillEventLog.HeaderID,
			event.LogFK:    LogKillEventLog.ID,
			// Oasis address
			constants.OfferId:    "604350",
			constants.PairColumn: "0xd257ccbe93e550a27236e8cc4971336f6cd2d53037ad567f10fbcc28df6a1eb1",
			// Maker address
			// PayGem address
			// BuyGem address
			constants.PayAmtColumn:    "2000000000000000000",
			constants.BuyAmtColumn:    "17999998000000",
			constants.TimestampColumn: "1583440730",
		},
	}
)

func LogKillModel() event.InsertionModel { return CopyModel(logKillModel) }
