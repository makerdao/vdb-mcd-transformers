package clip_yank

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type ClipYankEntity struct {
	Id              *big.Int
	ContractAddress common.Address
	HeaderID        int64
	LogID           int64
}
