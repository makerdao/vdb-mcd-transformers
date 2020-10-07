package test_data

import (
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
)

var rawNewCdpLog = types.Log{
	Address: common.HexToAddress(CdpManagerAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.NewCdpSignature()),
		common.HexToHash("0x000000000000000000000000a9fccb07dd3f774d5b9d02e99de1a27f47f91189"),
		common.HexToHash("0x000000000000000000000000a9fccb07dd3f774d5b9d02e99de1a27f47f91189"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000053"),
	},
	BlockNumber: uint64(int64(12975121)),
	TxHash:      common.HexToHash("0x4c2902029e9250a1927e096262bd6d23db0e0f3adef3a26cc4b3585e9ed86d52"),
	TxIndex:     112,
	BlockHash:   fakes.FakeHash,
	Index:       8,
	Removed:     false,
}

var NewCdpEventLog = core.EventLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawNewCdpLog,
	Transformed: false,
}

func NewCdpModel() event.InsertionModel { return CopyModel(newCdpModel) }

var newCdpModel = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.NewCdpTable,
	OrderedColumns: []event.ColumnName{
		event.HeaderFK, event.LogFK, constants.UsrColumn, constants.OwnColumn, constants.CdpColumn,
	},
	ColumnValues: event.ColumnValues{
		event.HeaderFK:      NewCdpEventLog.HeaderID,
		event.LogFK:         NewCdpEventLog.ID,
		constants.UsrColumn: "0xA9fCcB07DD3f774d5b9d02e99DE1a27f47F91189",
		constants.OwnColumn: "0xA9fCcB07DD3f774d5b9d02e99DE1a27f47F91189",
		constants.CdpColumn: "83",
	},
}
