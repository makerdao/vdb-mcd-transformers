package clip_take

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type ClipTakeEntity struct {
	Id              *big.Int
	Max             *big.Int
	Price           *big.Int
	Owe             *big.Int
	Tab             *big.Int
	Lot             *big.Int
	Usr             common.Address
	ContractAddress common.Address
	HeaderID        int64
	LogID           int64
}
