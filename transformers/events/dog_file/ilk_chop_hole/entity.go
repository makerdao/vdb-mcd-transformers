package ilk_chop_hole

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type DogFileIlkChopHoleEntity struct {
	Ilk  [32]byte
	What [32]byte
	Data *big.Int

	HeaderID        int64
	LogID           int64
	ContractAddress common.Address
}
