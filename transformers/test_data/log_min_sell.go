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
	logMinSellRawLog = types.Log{
		Address: common.HexToAddress(OasisAddresses()[0]),
		Topics: []common.Hash{
			common.HexToHash(constants.LogMinSellSignature()),
		},
		Data:        hexutil.MustDecode("0x0000000000000000000000006b175474e89094c44da98b954eedeac495271d0f0000000000000000000000000000000000000000000000001bc16d674ec80000"),
		BlockNumber: 8612271,
		TxHash:      common.HexToHash("0x367c11ca276af948c3170f1bf49cbd157163762735424116948244dcef76c058"),
		TxIndex:     156,
		BlockHash:   common.HexToHash("0x7d986f85263b1f48e196f35986857ffee532f8a43c15ecd32391a75f47b9ce70"),
		Index:       138,
		Removed:     false,
	}

	LogMinSellEventLog = core.EventLog{
		ID:          int64(rand.Int31()),
		HeaderID:    int64(rand.Int31()),
		Log:         logMinSellRawLog,
		Transformed: false,
	}

	LogMinSellPayGemAddress = common.HexToAddress("0x6b175474e89094c44da98b954eedeac495271d0f")

	logMinSellModel = event.InsertionModel{
		SchemaName: constants.MakerSchema,
		TableName:  constants.LogMinSellTable,
		OrderedColumns: []event.ColumnName{
			event.HeaderFK, event.LogFK, event.AddressFK, constants.PayGemColumn,
			constants.MinAmountColumn},
		ColumnValues: event.ColumnValues{
			event.HeaderFK: LogMinSellEventLog.HeaderID,
			event.LogFK:    LogMinSellEventLog.ID,
			// pay gem address id
			constants.MinAmountColumn: "2000000000000000000",
		},
	}
)

func LogMinSellModel() event.InsertionModel { return CopyModel(logMinSellModel) }
