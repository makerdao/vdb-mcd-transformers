package log_min_sell

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type LogMinSellEntity struct {
	PayGem          common.Address
	MinAmount       *big.Int
	HeaderID        int64
	LogID           int64
	ContractAddress common.Address
}
