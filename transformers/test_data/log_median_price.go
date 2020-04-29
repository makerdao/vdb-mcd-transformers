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
	rawEthLogMedianPriceLog = types.Log{
		Address: common.HexToAddress(EthMedianAddress()),
		Topics: []common.Hash{
			common.HexToHash(constants.LogMedianPriceSignature()),
		},
		Data:        hexutil.MustDecode("0x00000000000000000000000000000000000000000000000a708fc6189b4b8000000000000000000000000000000000000000000000000000000000005ea70222"),
		BlockNumber: 9955467,
		TxHash:      common.HexToHash("0xb4b8fb47e4423cde8548f4d2785839c34b89df1d07cb111bc670627ed915d2ad"),
		TxIndex:     0,
		BlockHash:   common.HexToHash("0xb24ada88d46a0d6a9d3c81744084b041b450afcfb48d16ba2475d4b8ea1ceaa0"),
		Index:       0,
		Removed:     false,
	}

	rawBatLogMedianPriceLog = types.Log{
		Address: common.HexToAddress((BatMedianAddress())),
		Topics: []common.Hash{
			common.HexToHash(constants.LogMedianPriceSignature()),
		},
		Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000027c16068847d100000000000000000000000000000000000000000000000000000000005ea7852e"),
		BlockNumber: 9957995,
		TxHash:      common.HexToHash("0x776152569b2e44d6b5d01b52082b82213df730cebcd03738f9387eb64fe760c8"),
		TxIndex:     0,
		BlockHash:   common.HexToHash("0xe618d33d78971a1cee37dde1d53c91cc0a86f21ae71b05cc95f9ee7f291efc23"),
		Index:       0,
		Removed:     false,
	}

	BatLogMedianPriceEventLog = core.EventLog{
		ID:          rand.Int63(),
		HeaderID:    rand.Int63(),
		Log:         rawBatLogMedianPriceLog,
		Transformed: false,
	}

	EthLogMedianPriceEventLog = core.EventLog{
		ID:          rand.Int63(),
		HeaderID:    rand.Int63(),
		Log:         rawEthLogMedianPriceLog,
		Transformed: false,
	}

	batLogMedianPriceModel = event.InsertionModel{
		SchemaName: constants.MakerSchema,
		TableName:  constants.LogMedianPriceTable,
		OrderedColumns: []event.ColumnName{
			event.HeaderFK,
			event.LogFK,
			constants.ValColumn,
			constants.AgeColumn,
		},
		ColumnValues: event.ColumnValues{
			event.HeaderFK:      BatLogMedianPriceEventLog.HeaderID,
			event.LogFK:         BatLogMedianPriceEventLog.ID,
			constants.ValColumn: "179042302500000000",
			constants.AgeColumn: "1588036910",
		},
	}

	ethLogMedianPriceModel = event.InsertionModel{
		SchemaName: constants.MakerSchema,
		TableName:  constants.LogMedianPriceTable,
		OrderedColumns: []event.ColumnName{
			event.HeaderFK,
			event.LogFK,
			constants.ValColumn,
			constants.AgeColumn,
		},
		ColumnValues: event.ColumnValues{
			event.HeaderFK:      EthLogMedianPriceEventLog.HeaderID,
			event.LogFK:         EthLogMedianPriceEventLog.ID,
			constants.ValColumn: "192578360000000000000",
			constants.AgeColumn: "1588003362",
		},
	}
)

func EthLogMedianPriceModel() event.InsertionModel { return CopyModel(ethLogMedianPriceModel) }
func BatLogMedianPriceModel() event.InsertionModel { return CopyModel(batLogMedianPriceModel) }
