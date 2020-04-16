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
	logUnsortedOfferRawLog = types.Log{
		Address:     common.HexToAddress(OasisAddresses()[0]),
		Topics:      []common.Hash{common.HexToHash(constants.LogUnsortedOfferSignature())},
		Data:        hexutil.MustDecode("0x00000000000000000000000000000000000000000000000000000000000af0fa"),
		BlockNumber: 9243052,
		TxHash:      common.HexToHash("0x8dff61bb629b1f45f31d8f91a644b8c3cdac48a3b85f04d3267f951a254d3968"),
		TxIndex:     157,
		BlockHash:   fakes.FakeHash,
		Index:       207,
		Removed:     false,
	}

	LogUnsortedOfferEventLog = core.EventLog{
		ID:          int64(rand.Int31()),
		HeaderID:    int64(rand.Int31()),
		Log:         logUnsortedOfferRawLog,
		Transformed: false,
	}

	logUnsortedOfferModel = event.InsertionModel{
		SchemaName: constants.MakerSchema,
		TableName:  constants.LogUnsortedOfferTable,
		OrderedColumns: []event.ColumnName{
			event.HeaderFK, event.LogFK, event.AddressFK, constants.OfferId,
		},
		ColumnValues: event.ColumnValues{
			event.HeaderFK:    LogUnsortedOfferEventLog.HeaderID,
			event.LogFK:       LogUnsortedOfferEventLog.ID,
			constants.OfferId: "717050",
		},
	}
)

func LogUnsortedOfferModel() event.InsertionModel { return CopyModel(logUnsortedOfferModel) }
