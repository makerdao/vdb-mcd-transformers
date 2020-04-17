package test_data

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"math/rand"
)

var (
	logBumpRawLog = types.Log{
		Address: common.HexToAddress(OasisAddresses()[0]),
		Topics: []common.Hash{
			common.HexToHash(constants.LogBumpSignature()),
			common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000045f23"),
			common.HexToHash("0x10aed75aa327f09ef87e5bdfaedf498ca260499a251ae5e049ddbd5e1633cd9c"),
			common.HexToHash("0x000000000000000000000000a4da0f347c6abe0e8bc71b5981fd92b364eda4c2"),
		},
		Data:        hexutil.MustDecode("0x00000000000000000000000089d24a6b4ccb1b6faa2625fe562bdd9a23260359000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc20000000000000000000000000000000000000000000000005563d8e666fa8000000000000000000000000000000000000000000000000000004a9b6384488000000000000000000000000000000000000000000000000000000000005d21091b"),
		BlockNumber: 8100012,
		TxHash:      common.HexToHash("0x0fcdfefc688424c83e79aaa238b71df27da55d5945920a6ea788b221e5c8bcf8"),
		TxIndex:     54,
		BlockHash:   common.HexToHash("0x4a21e04c8f5f74556e4eb6b5b89f805ff81761c2ade43663498278447fb52301"),
		Index:       31,
		Removed:     false,
	}

	LogBumpEventLog = core.EventLog{
		ID:          int64(rand.Int31()),
		HeaderID:    int64(rand.Int31()),
		Log:         logBumpRawLog,
		Transformed: false,
	}

	LogBumpPayGemAddress = common.HexToAddress("0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359")
	LogBumpBuyGemAddress = common.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2")

	logBumpModel = event.InsertionModel{
		SchemaName: constants.MakerSchema,
		TableName:  constants.LogBumpTable,
		OrderedColumns: []event.ColumnName{
			event.HeaderFK, event.LogFK, event.AddressFK, constants.OfferId, constants.PairColumn, constants.MakerColumn, constants.PayGemColumn,
			constants.BuyGemColumn, constants.PayAmtColumn, constants.BuyAmtColumn, constants.TimestampColumn,
		},
		ColumnValues: event.ColumnValues{
			event.HeaderFK: LogBumpEventLog.HeaderID,
			event.LogFK:    LogBumpEventLog.ID,
			// Oasis contract address id
			constants.OfferId:    "286499",
			constants.PairColumn: "0x10aed75aa327f09ef87e5bdfaedf498ca260499a251ae5e049ddbd5e1633cd9c",
			// Maker address id
			// Pay gem address id
			// Buy gem address id
			constants.PayAmtColumn:    "6153000000000000000",
			constants.BuyAmtColumn:    "21000000000000000",
			constants.TimestampColumn: "1562446107",
		},
	}
)

func LogBumpModel() event.InsertionModel { return CopyModel(logBumpModel) }
