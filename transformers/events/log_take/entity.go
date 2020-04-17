package log_take

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type LogTakeEntity struct {
	Id              common.Hash
	Pair            common.Hash
	Maker           common.Address
	PayGem          common.Address
	BuyGem          common.Address
	Taker           common.Address
	TakeAmt         *big.Int
	GiveAmt         *big.Int
	Timestamp       uint64
	HeaderID        int64
	LogID           int64
	ContractAddress common.Address
}
