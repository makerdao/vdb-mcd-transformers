package log_matching_enabled

import (
	"github.com/ethereum/go-ethereum/common"
)

type LogMatchingEnabledEntity struct {
	IsEnabled       bool
	HeaderID        int64
	LogID           int64
	ContractAddress common.Address
}
