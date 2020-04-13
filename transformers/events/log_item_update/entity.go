package log_item_update

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type LogItemUpdateEntity struct {
	Id              *big.Int
	HeaderID        int64
	LogID           int64
	ContractAddress common.Address
}
