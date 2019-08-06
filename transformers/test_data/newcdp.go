package test_data

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/mcd_transformers/transformers/events/new_cdp"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/big"
	"math/rand"
)

const (
	TemporaryNewCdpBlockNumber = int64(12975121)
	newCdpData                 = "0x0000000000000000000000000000000000000000000000000000000000000015"
	TemporaryNewCdpTransaction = "0x4c2902029e9250a1927e096262bd6d23db0e0f3adef3a26cc4b3585e9ed86d52"
)

var (
	newCdpUsr        = common.HexToAddress("0x000000000000000000000000a9fccb07dd3f774d5b9d02e99de1a27f47f91189")
	newCdpOwn        = common.HexToAddress("0x000000000000000000000000a9fccb07dd3f774d5b9d02e99de1a27f47f91189")
	newCdpCdp        = big.NewInt(21)
	newCdpRawJson, _ = json.Marshal(NewCdpHeaderSyncLog)
)

var rawNewCdpLog = types.Log{
	Address: common.HexToAddress(CdpManagerAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.NewCdpSignature()),
		common.HexToHash("0x000000000000000000000000a9fccb07dd3f774d5b9d02e99de1a27f47f91189"),
		common.HexToHash("0x000000000000000000000000a9fccb07dd3f774d5b9d02e99de1a27f47f91189"),
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

var NewCdpEntity = new_cdp.NewCdpEntity{
	Usr:      newCdpUsr,
	Own:      newCdpOwn,
	Cdp:      newCdpCdp,
	LogID:    NewCdpHeaderSyncLog.ID,
	HeaderID: NewCdpHeaderSyncLog.HeaderID,
}

var NewCdpModel = new_cdp.NewCdpModel{
	Usr:      "0xA9fCcB07DD3f774d5b9d02e99DE1a27f47F91189",
	Own:      "0xA9fCcB07DD3f774d5b9d02e99DE1a27f47F91189",
	Cdp:      newCdpCdp.String(),
	LogID:    NewCdpHeaderSyncLog.ID,
	HeaderID: NewCdpHeaderSyncLog.HeaderID,
}
