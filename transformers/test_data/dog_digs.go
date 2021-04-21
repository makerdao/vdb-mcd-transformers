package test_data

import (
	"math/big"
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
)

var (
	rad = big.NewInt(456)

	RawDogDigsLog = types.Log{
		Address: common.HexToAddress(Dog130Address()),
		Topics: []common.Hash{
			common.HexToHash(constants.DogDigsSignature()),
			common.HexToHash("0x4554482d41000000000000000000000000000000000000000000000000000000"), // ilk
		},
		Data: hexutil.MustDecode("0x" +
			"00000000000000000000000000000000000000000000000000000000000001c8" + //rad
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

	DogDigsEventLog = core.EventLog{
		ID:          int64(rand.Int31()),
		HeaderID:    int64(rand.Int31()),
		Log:         RawDogDigsLog,
		Transformed: false,
	}
	dogDigsModel = event.InsertionModel{
		SchemaName: constants.MakerSchema,
		TableName:  constants.DogDigsTable,
		OrderedColumns: []event.ColumnName{
			event.HeaderFK,
			event.LogFK,
			event.AddressFK,
			constants.IlkColumn,
			constants.RadColumn,
		},
		ColumnValues: event.ColumnValues{
			event.HeaderFK: DogDigsEventLog.HeaderID,
			event.LogFK:    DogDigsEventLog.ID,
			//event.AddressFK,
			//constants.IlkColumn,
			constants.RadColumn: rad.String(),
		},
	}
)

func DogDigsModel() event.InsertionModel { return CopyModel(dogDigsModel) }
