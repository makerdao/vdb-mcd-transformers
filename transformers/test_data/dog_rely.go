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
	RawDogRelyLog = types.Log{
		Address: common.HexToAddress(Dog130Address()),
		Topics: []common.Hash{
			common.HexToHash(constants.DogRelySignature()),
			common.HexToHash("0x000000000000000000000000dDb108893104dE4E1C6d0E47c42237dB4E617ACc"), //usr
		},
		Data:        hexutil.MustDecode("0x"),
		BlockNumber: testBlockNumber,
		TxHash:      common.Hash{},
		TxIndex:     0,
		BlockHash:   common.Hash{},
		Index:       0,
		Removed:     false,
	}

	DogRelyEventLog = core.EventLog{
		ID:          int64(rand.Int31()),
		HeaderID:    int64(rand.Int31()),
		Log:         RawDogDenyLog,
		Transformed: false,
	}
	dogRelyModel = event.InsertionModel{
		SchemaName: constants.MakerSchema,
		TableName:  constants.DogRelyTable,
		OrderedColumns: []event.ColumnName{
			event.HeaderFK,
			event.LogFK,
			event.AddressFK,
			constants.UsrColumn,
		},
		ColumnValues: event.ColumnValues{
			event.HeaderFK: DogRelyEventLog.HeaderID,
			event.LogFK:    DogRelyEventLog.ID,
			//event.AddressFK,
			//constants.UsrColumn,
		},
	}
)

func DogRelyModel() event.InsertionModel { return CopyModel(dogRelyModel) }
