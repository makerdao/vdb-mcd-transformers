package log_median_price

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type LogMedianPriceEntity struct {
	Val             *big.Int
	Age             *big.Int
	LogID           int64
	HeaderID        int64
	ContractAddress common.Address
}
