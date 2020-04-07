package log_item_update

import "math/big"

type LogItemUpdateEntity struct {
	Id       *big.Int
	HeaderID int64
	LogID    int64
}

