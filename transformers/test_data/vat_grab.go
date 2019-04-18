package test_data

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_grab"
)

var EthVatGrabLogWithPositiveDink = types.Log{
	Address: common.HexToAddress(KovanVatContractAddress),
	Topics: []common.Hash{
		common.HexToHash(KovanVatGrabSignature),
		common.HexToHash("0x66616b6520696c6b000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x00000000000000000000000007fa9ef6609ca7921112231f8f195138ebba2977"),
		common.HexToHash("0x0000000000000000000000007340e006f4135ba6970d43bf43d88dcad4e7a8ca"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e07bab3f4066616b6520696c6b00000000000000000000000000000000000000000000000000000000000000000000000007fa9ef6609ca7921112231f8f195138ebba29770000000000000000000000007340e006f4135ba6970d43bf43d88dcad4e7a8ca0000000000000000000000007526eb4f95e2a1394797cb38a921fb1eba09291b00000000000000000000000000000000000000000000003635c9adc5dea0000000000000000000000000000000000000000000000000006c6b935b8bbd40000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 94,
	TxHash:      common.HexToHash("0x18aa10dddd3543143d4eff314c8f3620c3b8299e57468d70ff7abd4498eee7db"),
	TxIndex:     4,
	BlockHash:   common.HexToHash("0x17b1de2797689d940a66911ebb9ae789528c269aea309c55128e38d48ddb37a9"),
	Index:       5,
	Removed:     false,
}

var rawVatGrabLogWithPositiveDink, _ = json.Marshal(EthVatGrabLogWithPositiveDink)
var VatGrabModelWithPositiveDink = vat_grab.VatGrabModel{
	Ilk:              "66616b6520696c6b000000000000000000000000000000000000000000000000",
	Urn:              "0x07Fa9eF6609cA7921112231F8f195138ebbA2977",
	V:                "0x7340e006f4135BA6970D43bf43d88DCAD4e7a8CA",
	W:                "0x7526EB4f95e2a1394797Cb38a921Fb1EbA09291B",
	Dink:             "1000000000000000000000",
	Dart:             "2000000000000000000000",
	LogIndex:         EthVatGrabLogWithPositiveDink.Index,
	TransactionIndex: EthVatGrabLogWithPositiveDink.TxIndex,
	Raw:              rawVatGrabLogWithPositiveDink,
}

var EthVatGrabLogWithNegativeDink = types.Log{
	Address: common.HexToAddress(KovanVatContractAddress),
	Topics: []common.Hash{
		common.HexToHash(KovanVatGrabSignature),
		common.HexToHash("0x66616b6520696c6b000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x00000000000000000000000007fa9ef6609ca7921112231f8f195138ebba2977"),
		common.HexToHash("0x0000000000000000000000007340e006f4135ba6970d43bf43d88dcad4e7a8ca"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e07bab3f4066616b6520696c6b00000000000000000000000000000000000000000000000000000000000000000000000007fa9ef6609ca7921112231f8f195138ebba29770000000000000000000000007340e006f4135ba6970d43bf43d88dcad4e7a8ca0000000000000000000000007526eb4f95e2a1394797cb38a921fb1eba09291bffffffffffffffffffffffffffffffffffffffffffffffc9ca36523a21600000ffffffffffffffffffffffffffffffffffffffffffffff93946ca47442c0000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 95,
	TxHash:      common.HexToHash("0x38c279812f842d7498ffe2f45f0d8c5de76ecd1ee7912636e7431129746b2e15"),
	TxIndex:     0,
	BlockHash:   common.HexToHash("0x4ea9f55ed4a97f686dd28a0d6d2dea8a0915a902a245e86516df3e4a57d1ca9d"),
	Index:       0,
	Removed:     false,
}

var rawVatGrabLogWithNegativeDink, _ = json.Marshal(EthVatGrabLogWithNegativeDink)
var VatGrabModelWithNegativeDink = vat_grab.VatGrabModel{
	Ilk:              "66616b6520696c6b000000000000000000000000000000000000000000000000",
	Urn:              "0x07Fa9eF6609cA7921112231F8f195138ebbA2977",
	V:                "0x7340e006f4135BA6970D43bf43d88DCAD4e7a8CA",
	W:                "0x7526EB4f95e2a1394797Cb38a921Fb1EbA09291B",
	Dink:             "-1000000000000000000000",
	Dart:             "-2000000000000000000000",
	LogIndex:         0,
	TransactionIndex: 0,
	Raw:              rawVatGrabLogWithNegativeDink,
}
