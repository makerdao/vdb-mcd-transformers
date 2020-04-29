package log_median_price

import "math/big"

type LogMedianPriceEntity struct {
	Val      *big.Int
	Age      *big.Int
	LogID    int64
	HeaderID int64
}
