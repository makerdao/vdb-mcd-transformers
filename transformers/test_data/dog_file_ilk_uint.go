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
	RawDogFileIlkUintLog = types.Log{
		Address: common.HexToAddress(Dog1xxAddress()),
		Topics: []common.Hash{
			common.HexToHash(constants.DogFileIlkUintSignature()),
			common.HexToHash("0x4554482d41000000000000000000000000000000000000000000000000000000"), // ilk
			common.HexToHash("0x00000000000000000000000000000000000000000000000000000000636c6970"), // what
		},
		Data: hexutil.MustDecode("0x" +
			"00000000000000000000000000000000000000000000000000000004FCF6BC30" +
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

	DogFileIlkUintEventLog = core.EventLog{
		ID:          int64(rand.Int31()),
		HeaderID:    int64(rand.Int31()),
		Log:         RawDogFileIlkUintLog,
		Transformed: false,
	}

	dogFileIlkUintModel = event.InsertionModel{
		SchemaName: constants.MakerSchema,
		TableName:  constants.DogFileIlkUintTable,
		OrderedColumns: []event.ColumnName{
			event.HeaderFK,
			event.LogFK,
			event.AddressFK,
			constants.IlkColumn,
			constants.WhatColumn,
			constants.DataColumn,
		},
		ColumnValues: event.ColumnValues{
			event.HeaderFK: DogFileIlkUintEventLog.HeaderID,
			event.LogFK:    DogFileIlkUintEventLog.ID,
			//event.AddressFK,
			//constants.IlkColumn
			constants.WhatColumn: DogFileIlkUintEventLog.Log.Topics[2].String(),
			constants.DataColumn: "21423897648",
		},
	}
)

func DogFileIlkUintModel() event.InsertionModel { return CopyModel(dogFileIlkUintModel) }
