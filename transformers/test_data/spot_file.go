package test_data

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/mcd_transformers/transformers/events/spot_file/pip"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

var EthSpotFilePipLog = types.Log{
	Address: common.HexToAddress(constants.SpotContractAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.SpotFilePipSignature()),
		common.HexToHash("0x0000000000000000000000004ba936a9338ae211300ea47899fbd111fd5dca31"),
		common.HexToHash("0x4554482d41000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x0000000000000000000000008c73ec0fbcdec6b8c060bc224d94740fd41f3774"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e0d4e8be834554482d410000000000000000000000000000000000000000000000000000000000000000000000000000008c73ec0fbcdec6b8c060bc224d94740fd41f3774000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 11257191,
	TxHash:      common.HexToHash("0xaae9e8bce346f86a01c5a3af137bc1f9bc7c0c767804a2b9b6356849aee0d7dd"),
	TxIndex:     1,
	BlockHash:   common.HexToHash("0xfa28e186578238fdd6b971add2ebe62a26dddf5ff971d50ee476c86b45362da1"),
	Index:       2,
	Removed:     false,
}

var rawSpotFilePipLog, _ = json.Marshal(EthSpotFilePipLog)
var SpotFilePipModel = pip.SpotFilePipModel{
	Ilk:              "4554482d41000000000000000000000000000000000000000000000000000000",
	Pip:              "0x8C73Ec0fBCdEC6b8C060BC224D94740FD41f3774",
	LogIndex:         EthSpotFilePipLog.Index,
	TransactionIndex: EthSpotFilePipLog.TxIndex,
	Raw:              rawSpotFilePipLog,
}
