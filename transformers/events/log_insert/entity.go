package log_insert

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type LogInsertEntity struct {
	Keeper          common.Address
	Id              *big.Int
	HeaderID        int64
	LogID           int64
	ContractAddress common.Address
}
