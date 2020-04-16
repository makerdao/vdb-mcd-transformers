package log_unsorted_offer

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type LogUnsortedOfferEntity struct {
	Id              *big.Int
	HeaderID        int64
	LogID           int64
	ContractAddress common.Address
}
