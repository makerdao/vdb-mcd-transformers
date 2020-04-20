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
	//TODO Updqte with real data
	logBuyEnabledRawLog = types.Log{
		Address:     common.HexToAddress(OasisAddresses()[0]),
		Topics:      []common.Hash{common.HexToHash(constants.LogBuyEnabledSignature())},
		Data:        hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000001"),
		BlockNumber: 9827070,
		TxHash:      common.HexToHash("0x6aef663b8483d1180faf1efee4501c27182c8496c1cb1615868af5cd324d2028"),
		TxIndex:     0,
		BlockHash:   fakes.FakeHash,
		Index:       0,
		Removed:     false,
	}

	LogBuyEnabledEventLog = core.EventLog{
		ID:          int64(rand.Int31()),
		HeaderID:    int64(rand.Int31()),
		Log:         logBuyEnabledRawLog,
		Transformed: false,
	}

	logBuyEnabledModel = event.InsertionModel{
		SchemaName: constants.MakerSchema,
		TableName:  constants.LogBuyEnabledTable,
		OrderedColumns: []event.ColumnName{
			event.HeaderFK, event.LogFK, event.AddressFK, constants.IsEnabled,
		},
		ColumnValues: event.ColumnValues{
			event.HeaderFK:      LogBuyEnabledEventLog.HeaderID,
			event.LogFK:         LogBuyEnabledEventLog.ID,
			constants.IsEnabled: true,
		},
	}
)

func LogBuyEnabledModel() event.InsertionModel { return CopyModel(logBuyEnabledModel) }
