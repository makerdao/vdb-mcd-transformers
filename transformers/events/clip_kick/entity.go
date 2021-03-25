package clip_kick

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type ClipKickEntity struct {
	Id              *big.Int
	Top             *big.Int
	Tab             *big.Int
	Lot             *big.Int
	Usr             common.Address
	Kpr             common.Address
	Coin            *big.Int
	ContractAddress common.Address
	HeaderID        int64
	LogID           int64
}
