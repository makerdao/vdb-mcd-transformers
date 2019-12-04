package test_data

import (
	"encoding/json"
	"math/big"
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
)

const (
	TemporaryNewCdpBlockNumber = int64(12975121)
	newCdpData                 = "0x0000000000000000000000000000000000000000000000000000000000000015"
	TemporaryNewCdpTransaction = "0x4c2902029e9250a1927e096262bd6d23db0e0f3adef3a26cc4b3585e9ed86d52"
)

var (
	newCdpUsr        = common.HexToAddress("0x000000000000000000000000a9fccb07dd3f774d5b9d02e99de1a27f47f91189")
	newCdpOwn        = common.HexToAddress("0x000000000000000000000000a9fccb07dd3f774d5b9d02e99de1a27f47f91189")
	newCdpCdp        = big.NewInt(83)
	newCdpRawJson, _ = json.Marshal(NewCdpHeaderSyncLog)
)

var rawNewCdpLog = types.Log{
	Address: common.HexToAddress(CdpManagerAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.NewCdpSignature()),
		common.HexToHash("0x000000000000000000000000a9fccb07dd3f774d5b9d02e99de1a27f47f91189"),
		common.HexToHash("0x000000000000000000000000a9fccb07dd3f774d5b9d02e99de1a27f47f91189"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000053"),
	},
	Data:        hexutil.MustDecode(newCdpData),
	BlockNumber: uint64(TemporaryNewCdpBlockNumber),
	TxHash:      common.HexToHash(TemporaryNewCdpTransaction),
	TxIndex:     112,
	BlockHash:   fakes.FakeHash,
	Index:       8,
	Removed:     false,
}

var NewCdpHeaderSyncLog = core.HeaderSyncLog{
	ID:          int64(rand.Int31()),
	HeaderID:    int64(rand.Int31()),
	Log:         rawNewCdpLog,
	Transformed: false,
}

func NewCdpModel() shared.InsertionModel { return CopyModel(newCdpModel) }

var newCdpModel = shared.InsertionModel{
	SchemaName: constants.MakerSchema,
	TableName:  "new_cdp",
	OrderedColumns: []string{
		constants.HeaderFK, constants.LogFK, "usr", "own", "cdp",
	},
	ColumnValues: shared.ColumnValues{
		constants.HeaderFK: NewCdpHeaderSyncLog.HeaderID,
		constants.LogFK:    NewCdpHeaderSyncLog.ID,
		"usr":              "0xA9fCcB07DD3f774d5b9d02e99DE1a27f47F91189",
		"own":              "0xA9fCcB07DD3f774d5b9d02e99DE1a27f47F91189",
		"cdp":              newCdpCdp.String(),
	},
	ForeignKeyValues: shared.ForeignKeyValues{},
}
