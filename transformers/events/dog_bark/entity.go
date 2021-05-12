package dog_bark

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type DogBarkEntity struct {
	Ilk  [32]byte
	Urn  common.Address
	Ink  *big.Int
	Art  *big.Int
	Due  *big.Int
	Clip common.Address
	Id   *big.Int

	HeaderID        int64
	LogID           int64
	ContractAddress common.Address
}
