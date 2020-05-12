package shared

import (
	"errors"
	"math/big"
)

func IsZero(n *big.Int) bool {
	zero := big.NewInt(0)
	return n.Cmp(zero) == 0
}

func StringToBigInt(s string) (*big.Int, error) {
	n, ok := big.NewInt(0).SetString(s, 10)
	if !ok {
		return nil, errors.New("error formatting string as *big.Int")
	}
	return n, nil
}
