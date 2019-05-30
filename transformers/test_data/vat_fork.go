package test_data

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/mcd_transformers/transformers/events/vat_fork"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

var EthVatForkLogWithNegativeDart = types.Log{
	Address: common.HexToAddress(constants.VatContractAddress()),
	Topics: []common.Hash{
		common.HexToHash("0x00000000000000000000000000000000000000000000000000000000870c616d"),
		common.HexToHash("0x66616b6520696c6b000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x00000000000000000000000007fa9ef6609ca7921112231f8f195138ebba2977"),
		common.HexToHash("0x0000000000000000000000007526eb4f95e2a1394797cb38a921fb1eba09291b"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e0870c616d66616b6520696c6b00000000000000000000000000000000000000000000000000000000000000000000000007fa9ef6609ca7921112231f8f195138ebba29770000000000000000000000007526eb4f95e2a1394797cb38a921fb1eba09291bffffffffffffffffffffffffffffffffffffffffffffffc9ca36523a21600000ffffffffffffffffffffffffffffffffffffffffffffff93946ca47442c00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 146,
	TxHash:      common.HexToHash("0xe64bdd39a752e1911e841d634a6fa8d4ef229a03f0555f9e055caec1ae4930c2"),
	TxIndex:     1,
	BlockHash:   common.HexToHash("0xf31c6d2dadd23f408e5158dce47ba20fef8c17bc60af6e1f35a89769bc20d6f0"),
	Index:       2,
	Removed:     false,
}

var rawVatForkLogWithNegativeDart, _ = json.Marshal(EthVatForkLogWithNegativeDart)
var VatForkModelWithNegativeDart = vat_fork.VatForkModel{
	Ilk:              "66616b6520696c6b000000000000000000000000000000000000000000000000",
	Src:              "0x07Fa9eF6609cA7921112231F8f195138ebbA2977",
	Dst:              "0x7526EB4f95e2a1394797Cb38a921Fb1EbA09291B",
	Dink:             "-1000000000000000000000",
	Dart:             "-2000000000000000000000",
	LogIndex:         EthVatForkLogWithNegativeDart.Index,
	TransactionIndex: EthVatForkLogWithNegativeDart.TxIndex,
	Raw:              rawVatForkLogWithNegativeDart,
}
