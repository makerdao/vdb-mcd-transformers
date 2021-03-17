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
	RawDogFileIlkChopHoleLog = types.Log{
		Address: common.HexToAddress(Dog1xxAddress()),
		Topics: []common.Hash{
			common.HexToHash(constants.DogFileIlkChopHoleSignature()),
			common.HexToHash("0x4554482d41000000000000000000000000000000000000000000000000000000"), // ilk
			common.HexToHash("0x00000000000000000000000000000000000000000000000000000000636c6970"), // what
		},
		Data: hexutil.MustDecode("0x" +
			"00000000000000000000000000000000000000000000000000000004FCF6BC30"),
		BlockNumber: testBlockNumber,
		TxHash:      common.Hash{},
		TxIndex:     0,
		BlockHash:   common.Hash{},
		Index:       0,
		Removed:     false,
	}

	DogFileIlkChopHoleEventLog = core.EventLog{
		ID:          int64(rand.Int31()),
		HeaderID:    int64(rand.Int31()),
		Log:         RawDogFileIlkChopHoleLog,
		Transformed: false,
	}

	dogFileIlkChopHoleModel = event.InsertionModel{
		SchemaName: constants.MakerSchema,
		TableName:  constants.DogFileIlkChopHoleTable,
		OrderedColumns: []event.ColumnName{
			event.HeaderFK,
			event.LogFK,
			event.AddressFK,
			constants.IlkColumn,
			constants.WhatColumn,
			constants.DataColumn,
		},
		ColumnValues: event.ColumnValues{
			event.HeaderFK: DogFileIlkChopHoleEventLog.HeaderID,
			event.LogFK:    DogFileIlkChopHoleEventLog.ID,
			//event.AddressFK,
			//constants.IlkColumn
			constants.WhatColumn: DogFileIlkChopHoleEventLog.Log.Topics[2].String(),
			constants.DataColumn: "21423897648",
		},
	}
)

func DogFileIlkChopHoleModel() event.InsertionModel { return CopyModel(dogFileIlkChopHoleModel) }
