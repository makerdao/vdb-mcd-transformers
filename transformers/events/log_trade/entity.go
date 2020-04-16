package log_trade

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type LogTradeEntity struct {
	PayGem          common.Address
	BuyGem          common.Address
	PayAmt          *big.Int
	BuyAmt          *big.Int
	HeaderID        int64
	LogID           int64
	ContractAddress common.Address
}
