package log_buy_enabled

import (
	"github.com/ethereum/go-ethereum/common"
)

type LogBuyEnabledEntity struct {
	IsEnabled       bool
	HeaderID        int64
	LogID           int64
	ContractAddress common.Address
}
