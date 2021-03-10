package ilk_uint

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type DogFileIlkUintEntity struct {
	Ilk  [32]byte
	What [32]byte
	Data *big.Int

	HeaderID        int64
	LogID           int64
	ContractAddress common.Address
}
