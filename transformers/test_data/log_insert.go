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
	logInsertRawLog = types.Log{
		Address: common.HexToAddress(OasisAddresses()[0]),
		Topics: []common.Hash{
			common.HexToHash(constants.LogInsertSignature()),
		},
		Data:        hexutil.MustDecode("0x0000000000000000000000003a32292c53bf42b6317334392bf0272da29832520000000000000000000000000000000000000000000000000000000000008d92"),
		BlockNumber: 7352952,
		TxHash:      common.HexToHash("0xd42e7dbf03aff8bdca2fbfe3c0be6043c7ab00df3401b4d767d0d13cb28fa7d8"),
		TxIndex:     44,
		BlockHash:   common.HexToHash("0xab25258d75f36031c6ac3153292a0636ccb580b6ad65269588ce390c692b7311"),
		Index:       51,
		Removed:     false,
	}

	LogInsertEventLog = core.EventLog{
		ID:          int64(rand.Int31()),
		HeaderID:    int64(rand.Int31()),
		Log:         logInsertRawLog,
		Transformed: false,
	}

	LogInsertKeeperAddress = common.HexToAddress("0x3a32292c53bf42b6317334392bf0272da2983252")

	logInsertModel = event.InsertionModel{
		SchemaName: constants.MakerSchema,
		TableName:  constants.LogInsertTable,
		OrderedColumns: []event.ColumnName{
			event.HeaderFK, event.LogFK, event.AddressFK, constants.KeeperColumn,
			constants.OfferId},
		ColumnValues: event.ColumnValues{
			event.HeaderFK: LogInsertEventLog.HeaderID,
			event.LogFK:    LogInsertEventLog.ID,
			// Keeper address id
			constants.OfferId: "36242",
		},
	}
)

func LogInsertModel() event.InsertionModel { return CopyModel(logInsertModel) }
