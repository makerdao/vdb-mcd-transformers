package test_data

import (
	"math/big"
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/log_item_update"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
)

var (
	logItemUpdateOfferId, _ = new(big.Int).SetString("228696", 10)
	logItemUpdateRawLog     = types.Log{
		Address:     common.HexToAddress(OasisAddresses()[0]),
		Topics:      []common.Hash{common.HexToHash(constants.LogItemUpdateSignature())},
		Data:        hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000037d58"),
		BlockNumber: 9827070,
		TxHash:      common.HexToHash("0x6aef663b8483d1180faf1efee4501c27182c8496c1cb1615868af5cd324d2028"),
		TxIndex:     0,
		BlockHash:   fakes.FakeHash,
		Index:       0,
		Removed:     false,
	}

	LogItemUpdateEventLog = core.EventLog{
		ID:          int64(rand.Int31()),
		HeaderID:    int64(rand.Int31()),
		Log:         logItemUpdateRawLog,
		Transformed: false,
	}

	logItemUpdateModel = event.InsertionModel{
		SchemaName: constants.MakerSchema,
		TableName:  constants.LogItemUpdateTable,
		OrderedColumns: []event.ColumnName{
			event.HeaderFK, event.LogFK, event.AddressFK, log_item_update.OfferId,
		},
		ColumnValues: event.ColumnValues{
			event.HeaderFK:          LogItemUpdateEventLog.HeaderID,
			event.LogFK:             LogItemUpdateEventLog.ID,
			log_item_update.OfferId: logItemUpdateOfferId.String(),
		},
	}
)

func LogItemUpdateModel() event.InsertionModel { return CopyModel(logItemUpdateModel) }
