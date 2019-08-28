package new_cdp

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type NewCdpEntity struct {
	Usr              common.Address
	Own              common.Address
	Cdp              *big.Int
	LogIndex         uint
	TransactionIndex uint
	Raw              types.Log
}
