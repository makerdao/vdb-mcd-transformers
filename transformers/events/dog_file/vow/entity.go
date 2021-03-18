package vow

import (
	"github.com/ethereum/go-ethereum/common"
)

type DogFileVowEntity struct {
	What [32]byte
	Data common.Address

	HeaderID        int64
	LogID           int64
	ContractAddress common.Address
}
