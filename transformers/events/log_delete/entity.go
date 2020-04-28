package log_delete

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type LogDeleteEntity struct {
	Keeper          common.Address
	Id              *big.Int
	HeaderID        int64
	LogID           int64
	ContractAddress common.Address
}
