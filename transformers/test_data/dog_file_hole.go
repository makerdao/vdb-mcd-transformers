package test_data

import (
	"math/big"
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
)

var holeData = big.NewInt(1234)

var (
	RawDogFileHoleLog = types.Log{
		Address: common.HexToAddress(Dog1xxAddress()),
		Topics: []common.Hash{
			common.HexToHash(constants.DogFileHoleSignature()),
			common.HexToHash("0x486f6c6500000000000000000000000000000000000000000000000000000000"), // what
		},
		Data: hexutil.MustDecode("0x" +
			"00000000000000000000000000000000000000000000000000000000000004D2" + //data
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

	DogFileHoleEventLog = core.EventLog{
		ID:          int64(rand.Int31()),
		HeaderID:    int64(rand.Int31()),
		Log:         RawDogFileHoleLog,
		Transformed: false,
	}

	dogFileHoleModel = event.InsertionModel{
		SchemaName: constants.MakerSchema,
		TableName:  constants.DogFileHoleTable,
		OrderedColumns: []event.ColumnName{
			event.HeaderFK,
			event.LogFK,
			event.AddressFK,
			constants.WhatColumn,
			constants.DataColumn,
		},
		ColumnValues: event.ColumnValues{
			event.HeaderFK: DogFileHoleEventLog.HeaderID,
			event.LogFK:    DogFileHoleEventLog.ID,
			//event.AddressFK,
			constants.WhatColumn: shared.DecodeHexToText(DogFileHoleEventLog.Log.Topics[1].Hex()),
			constants.DataColumn: holeData.String(),
		},
	}
)

func DogFileHoleModel() event.InsertionModel { return CopyModel(dogFileHoleModel) }
