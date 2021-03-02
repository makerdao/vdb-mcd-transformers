package ilk_clip

import (
	"github.com/ethereum/go-ethereum/common"
)

type DogFileIlkClipEntity struct {
	Ilk  [32]byte
	What [32]byte
	Clip common.Address

	HeaderID        int64
	LogID           int64
	ContractAddress common.Address
}
