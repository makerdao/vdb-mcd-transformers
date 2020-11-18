package log_make

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type LogMakeEntity struct {
	Id              [32]byte
	Pair            common.Hash
	Maker           common.Address
	PayGem          common.Address
	BuyGem          common.Address
	PayAmt          *big.Int
	BuyAmt          *big.Int
	Timestamp       uint64
	HeaderID        int64
	LogID           int64
	ContractAddress common.Address
}
