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
	logTakeRawLog = types.Log{
		Address: common.HexToAddress(OasisAddresses()[0]),
		Topics: []common.Hash{
			common.HexToHash(constants.LogMakeSignature()),
			common.HexToHash("0xcdd6659bca20e2b28ea10ead902280762ac8977c84459a152f90e561d50edf8c"),
			common.HexToHash("0x0000000000000000000000006ff7d252627d35b8eb02607c8f27acdb18032718"),
			common.HexToHash("0x0000000000000000000000003a32292c53bf42b6317334392bf0272da2983252"),
		},
		Data:        hexutil.MustDecode("0x00000000000000000000000000000000000000000000000000000000000c626c000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc20000000000000000000000006b175474e89094c44da98b954eedeac495271d0f00000000000000000000000000000000000000000000000039bb49f599a000000000000000000000000000000000000000000000000000319aa46324eae00000000000000000000000000000000000000000000000000000000000005e3e3b4f"),
		BlockNumber: 9439987,
		TxHash:      common.HexToHash("0x18883e113519d3f14d9d5c8e3c0ca84dea77fa8405443f3e86652ce96ddd20af"),
		TxIndex:     166,
		BlockHash:   common.HexToHash("0xa7a45cdbc69c1e8ae99226483817aed6fc8022cb69daac15fafff505d1249b3b"),
		Index:       125,
		Removed:     false,
	}

	LogTakeEventLog = core.EventLog{
		ID:          int64(rand.Int31()),
		HeaderID:    int64(rand.Int31()),
		Log:         logTakeRawLog,
		Transformed: false,
	}

	LogTakePayGemAddress = common.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2")
	LogTakeBuyGemAddress = common.HexToAddress("0x6b175474e89094c44da98b954eedeac495271d0f")

	logTakeModel = event.InsertionModel{
		SchemaName: constants.MakerSchema,
		TableName:  constants.LogTakeTable,
		OrderedColumns: []event.ColumnName{
			event.HeaderFK, event.LogFK, event.AddressFK, constants.OfferId, constants.PairColumn,
			constants.MakerColumn, constants.PayGemColumn, constants.BuyGemColumn, constants.TakerColumn,
			constants.TakeAmtColumn, constants.GiveAmtColumn, constants.TimestampColumn,
		},
		ColumnValues: event.ColumnValues{
			event.HeaderFK: LogTakeEventLog.HeaderID,
			event.LogFK:    LogTakeEventLog.ID,
			// Oasis contract address id
			constants.OfferId:    "811628",
			constants.PairColumn: "0xcdd6659bca20e2b28ea10ead902280762ac8977c84459a152f90e561d50edf8c",
			// Maker address id
			// Pay gem address id
			// Buy gem address id
			// Taker address id
			constants.TakeAmtColumn:   "4160000000000000000",
			constants.GiveAmtColumn:   "915033600000000000000",
			constants.TimestampColumn: "1581136719",
		},
	}
)

func LogTakeModel() event.InsertionModel { return CopyModel(logTakeModel) }
