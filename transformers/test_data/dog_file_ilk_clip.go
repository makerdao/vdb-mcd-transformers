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
	ClipAddress     = "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
	testBlockNumber = uint64(4)

	RawDogFileIlkClipLog = types.Log{
		Address: common.HexToAddress(Dog1xxAddress()),
		Topics: []common.Hash{
			common.HexToHash(constants.DogFileIlkClipSignature()),
			common.HexToHash("0x4554482d41000000000000000000000000000000000000000000000000000000"), // ilk
			common.HexToHash("0x00000000000000000000000007fa9ef6609ca7921112231f8f195138ebba2977"), // what
			common.HexToHash("0x000000000000000000000000be8e3e3618f7474f8cb1d074a26affef007e98fb"), // clip
		},
		Data: hexutil.MustDecode("0x" +
			"0000000000000000000000000000000000000000000000000000000000000000" +
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

	DogFileIlkClipEventLog = core.EventLog{
		ID:          int64(rand.Int31()),
		HeaderID:    int64(rand.Int31()),
		Log:         RawDogFileIlkClipLog,
		Transformed: false,
	}

	dogFileIlkClipModel = event.InsertionModel{
		SchemaName: constants.MakerSchema,
		TableName:  constants.DogFileIlkClipTable,
		OrderedColumns: []event.ColumnName{
			event.HeaderFK,
			event.LogFK,
			event.AddressFK,
			constants.IlkColumn,
			constants.WhatColumn,
			constants.ClipColumn,
		},
		ColumnValues: event.ColumnValues{
			event.HeaderFK: DogFileIlkClipEventLog.HeaderID,
			event.LogFK:    DogFileIlkClipEventLog.ID,
			//event.AddressFK,
			//constants.IlkColumn
			constants.WhatColumn: DogFileIlkClipEventLog.Log.Topics[2],
			//constants.ClipColumn,
		},
	}
)

func DogFileIlkClipModel() event.InsertionModel { return CopyModel(dogFileIlkClipModel) }
