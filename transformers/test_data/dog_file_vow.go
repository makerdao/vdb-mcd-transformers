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
	DataAddress      = "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
	RawDogFileVowLog = types.Log{
		Address: common.HexToAddress(Dog1xxAddress()),
		Topics: []common.Hash{
			common.HexToHash(constants.DogFileVowSignature()),
			common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000766f77"), // what
		},
		Data: hexutil.MustDecode("0x" +
			"000000000000000000000000BE8E3e3618f7474F8cB1d074A26afFef007E98FB" +
			"0000000000000000000000000000000000000000000000000000000000000000" +
			"0000000000000000000000000000000000000000000000000000000000000000" +
			"0000000000000000000000000000000000000000000000000000000000000000"),
		BlockNumber: testBlockNumber,
		TxHash:      common.Hash{},
		TxIndex:     0,
		BlockHash:   common.Hash{},
		Index:       0,
		Removed:     false,
	}

	DogFileVowEventLog = core.EventLog{
		ID:          int64(rand.Int31()),
		HeaderID:    int64(rand.Int31()),
		Log:         RawDogFileVowLog,
		Transformed: false,
	}

	dogFileVowModel = event.InsertionModel{
		SchemaName: constants.MakerSchema,
		TableName:  constants.DogFileVowTable,
		OrderedColumns: []event.ColumnName{
			event.HeaderFK,
			event.LogFK,
			event.AddressFK,
			constants.WhatColumn,
			constants.DataColumn,
		},
		ColumnValues: event.ColumnValues{
			event.HeaderFK: DogFileVowEventLog.HeaderID,
			event.LogFK:    DogFileVowEventLog.ID,
			//event.AddressFK,
			constants.WhatColumn: DogFileVowEventLog.Log.Topics[1].String(),
			//constants.DataColumn,
		},
	}
)

func DogFileVowModel() event.InsertionModel { return CopyModel(dogFileVowModel) }
