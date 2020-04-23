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
	logMakeRawLog = types.Log{
		Address: common.HexToAddress(OasisAddresses()[0]),
		Topics: []common.Hash{
			common.HexToHash(constants.LogMakeSignature()),
			common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000c6280"),
			common.HexToHash("0xcdd6659bca20e2b28ea10ead902280762ac8977c84459a152f90e561d50edf8c"),
			common.HexToHash("0x0000000000000000000000006ff7d252627d35b8eb02607c8f27acdb18032718"),
		},
		Data:        hexutil.MustDecode("0x000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc20000000000000000000000006b175474e89094c44da98b954eedeac495271d0f0000000000000000000000000000000000000000000000037d73b4e66fa60000000000000000000000000000000000000000000000000304f0cf4f5072f08000000000000000000000000000000000000000000000000000000000005e3e55fc"),
		BlockNumber: 9440502,
		TxHash:      common.HexToHash("0xb2f2f13a6ef0d1d6dd153d1a019929a31799498dcb1c100e04695fc0e95b9e58"),
		TxIndex:     1772,
		BlockHash:   common.HexToHash("0xe899b35b2b4fa50eb67a7d302978ba52338e8502dfbb308b0deeb49f11b74b09"),
		Index:       89,
		Removed:     false,
	}

	LogMakeEventLog = core.EventLog{
		ID:          int64(rand.Int31()),
		HeaderID:    int64(rand.Int31()),
		Log:         logMakeRawLog,
		Transformed: false,
	}

	LogMakePayGemAddress = common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2")
	LogMakeBuyGemAddress = common.HexToAddress("0x6B175474E89094C44Da98b954EedeAC495271d0F")

	logMakeModel = event.InsertionModel{
		SchemaName: constants.MakerSchema,
		TableName:  constants.LogMakeTable,
		OrderedColumns: []event.ColumnName{
			event.HeaderFK, event.LogFK, event.AddressFK, constants.OfferId, constants.PairColumn, constants.MakerColumn, constants.PayGemColumn,
			constants.BuyGemColumn, constants.PayAmtColumn, constants.BuyAmtColumn, constants.TimestampColumn,
		},
		ColumnValues: event.ColumnValues{
			event.HeaderFK: LogMakeEventLog.HeaderID,
			event.LogFK:    LogMakeEventLog.ID,
			// Oasis contract address id
			constants.OfferId:    "811648",
			constants.PairColumn: "0xcdd6659bca20e2b28ea10ead902280762ac8977c84459a152f90e561d50edf8c",
			// Maker address id
			// Pay gem address id
			// Buy gem address id
			constants.PayAmtColumn:    "64380000000000000000",
			constants.BuyAmtColumn:    "14258238600000000000000",
			constants.TimestampColumn: "1581143548",
		},
	}
)

func LogMakeModel() event.InsertionModel { return CopyModel(logMakeModel) }
