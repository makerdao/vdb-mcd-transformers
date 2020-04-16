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
	logTradeRawLog = types.Log{
		Address: common.HexToAddress(OasisAddresses()[0]),
		Topics: []common.Hash{
			common.HexToHash(constants.LogTradeSignature()),
			common.HexToHash("0x000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"),
			common.HexToHash("0x0000000000000000000000006b175474e89094c44da98b954eedeac495271d0f"),
		},
		Data:        hexutil.MustDecode("0x0000000000000000000000000000000000000000000000003e96bd1ae10300000000000000000000000000000000000000000000000000350dc64bc8bb8b0000"),
		BlockNumber: 9439704,
		TxHash:      common.HexToHash("0x5cb3965e1c4c65692e25b98832fd2957125d9062573c9e3bd454bb154b4e31bc"),
		TxIndex:     179,
		BlockHash:   common.HexToHash("0x1f77666bf62aff76f5abfd635682b4de8eb49879187ec075b46a8f33e59991b9"),
		Index:       168,
		Removed:     false,
	}

	LogTradeEventLog = core.EventLog{
		ID:          int64(rand.Int31()),
		HeaderID:    int64(rand.Int31()),
		Log:         logTradeRawLog,
		Transformed: false,
	}

	logTradeModel = event.InsertionModel{
		SchemaName: constants.MakerSchema,
		TableName:  constants.LogTradeTable,
		OrderedColumns: []event.ColumnName{
			event.HeaderFK, event.LogFK, event.AddressFK, constants.PayGemColumn, constants.BuyGemColumn,
			constants.PayAmtColumn, constants.BuyAmtColumn,
		},
		ColumnValues: event.ColumnValues{
			event.HeaderFK: LogTradeEventLog.HeaderID,
			event.LogFK:    LogTradeEventLog.ID,
			// Oasis contract address id
			// Pay gem address id
			// Buy gem address id
			constants.PayAmtColumn: "4510000000000000000",
			constants.BuyAmtColumn: "978670000000000000000",
		},
	}
)

func LogTradeModel() event.InsertionModel { return CopyModel(logTradeModel) }
