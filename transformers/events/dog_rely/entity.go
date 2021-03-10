package dog_rely

import (
	"github.com/ethereum/go-ethereum/common"
)

type DogRelyEntity struct {
	Usr common.Address

	HeaderID        int64
	LogID           int64
	ContractAddress common.Address
}
