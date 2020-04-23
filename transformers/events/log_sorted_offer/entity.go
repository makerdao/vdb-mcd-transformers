package log_sorted_offer

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type LogSortedOfferEntity struct {
	Id              *big.Int
	HeaderID        int64
	LogID           int64
	ContractAddress common.Address
}
