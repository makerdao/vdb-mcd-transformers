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
	logSortedOfferRawLog = types.Log{
		Address:     common.HexToAddress(OasisAddresses()[0]),
		Topics:      []common.Hash{common.HexToHash(constants.LogSortedOfferSignature())},
		Data:        hexutil.MustDecode("0x00000000000000000000000000000000000000000000000000000000000c6280"),
		BlockNumber: 9440502,
		TxHash:      common.HexToHash("0xb2f2f13a6ef0d1d6dd153d1a019929a31799498dcb1c100e04695fc0e95b9e58"),
		TxIndex:     124,
		BlockHash:   fakes.FakeHash,
		Index:       97,
		Removed:     false,
	}

	LogSortedOfferEventLog = core.EventLog{
		ID:          int64(rand.Int31()),
		HeaderID:    int64(rand.Int31()),
		Log:         logSortedOfferRawLog,
		Transformed: false,
	}

	logSortedOfferModel = event.InsertionModel{
		SchemaName: constants.MakerSchema,
		TableName:  constants.LogSortedOfferTable,
		OrderedColumns: []event.ColumnName{
			event.HeaderFK, event.LogFK, event.AddressFK, constants.OfferId,
		},
		ColumnValues: event.ColumnValues{
			event.HeaderFK:    LogSortedOfferEventLog.HeaderID,
			event.LogFK:       LogSortedOfferEventLog.ID,
			constants.OfferId: "811648",
		},
	}
)

func LogSortedOfferModel() event.InsertionModel { return CopyModel(logSortedOfferModel) }
