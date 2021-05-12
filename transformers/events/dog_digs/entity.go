package dog_digs

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type DogDigsEntity struct {
	Ilk [32]byte
	Rad *big.Int

	HeaderID        int64
	LogID           int64
	ContractAddress common.Address
}
