package new_cdp

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type NewCdpEntity struct {
	Usr      common.Address
	Own      common.Address
	Cdp      *big.Int
	HeaderID int64
	LogID    int64
}
