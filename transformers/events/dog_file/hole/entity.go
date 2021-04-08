package hole

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type DogFileHoleEntity struct {
	What [32]byte
	Data *big.Int

	HeaderID        int64
	LogID           int64
	ContractAddress common.Address
}
