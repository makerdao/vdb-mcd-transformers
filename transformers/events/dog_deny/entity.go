package dog_deny

import (
	"github.com/ethereum/go-ethereum/common"
)

type DogDenyEntity struct {
	Usr common.Address

	HeaderID        int64
	LogID           int64
	ContractAddress common.Address
}
