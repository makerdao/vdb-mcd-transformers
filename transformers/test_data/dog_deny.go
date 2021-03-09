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
	RawDogDenyLog = types.Log{
		Address: common.HexToAddress(Dog1xxAddress()),
		Topics: []common.Hash{
			common.HexToHash(constants.DogDenySignature()),
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

	DogDenyEventLog = core.EventLog{
		ID:          int64(rand.Int31()),
		HeaderID:    int64(rand.Int31()),
		Log:         RawDogDenyLog,
		Transformed: false,
	}
	dogDenyModel = event.InsertionModel{
		SchemaName: constants.MakerSchema,
		TableName:  constants.DogDenyTable,
		OrderedColumns: []event.ColumnName{
			event.HeaderFK,
			event.LogFK,
			event.AddressFK,
			constants.UsrColumn,
		},
		ColumnValues: event.ColumnValues{
			event.HeaderFK: DogDenyEventLog.HeaderID,
			event.LogFK:    DogDenyEventLog.ID,
			//event.AddressFK,
			//constants.UsrColumn,
		},
	}
)

func DogDenyModel() event.InsertionModel { return CopyModel(dogDenyModel) }
