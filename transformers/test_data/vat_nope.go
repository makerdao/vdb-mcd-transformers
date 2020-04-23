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

var rawVatNopeLog = types.Log{
	Address: common.HexToAddress(VatAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.VatHopeSignature()),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e0dc4d20fa00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 8928187,
	TxHash:      common.HexToHash("0x06cf0716136996c21e79d7cd47c5a534bc680e8828a5ff703c6eac74c267510e"),
	TxIndex:     1238,
	BlockHash:   common.HexToHash("0xce33eb7e3c9e0206e6d57dcd03785fe08d4108f94a806a946bf90cb22d1d81ff"),
	Index:       65,
	Removed:     false,
}

var VatNopeEventLog = core.EventLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawVatNopeLog,
	Transformed: false,
}

var vatNopeModel = event.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  constants.VatNopeTable,
	OrderedColumns: []event.ColumnName{
		event.HeaderFK, event.LogFK, constants.UsrColumn,
	},
	ColumnValues: event.ColumnValues{
		event.HeaderFK: VatNopeEventLog.HeaderID,
		event.LogFK:    VatNopeEventLog.ID,
		//constants.UsrColumn
	},
}

func VatNopeModel() event.InsertionModel { return CopyModel(vatNopeModel) }
