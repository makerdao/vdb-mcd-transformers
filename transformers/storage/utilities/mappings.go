package utilities

import (
	"errors"
	"fmt"
)

var ErrInvalidAddress = func(addr string) error {
	return errors.New(fmt.Sprintf("invalid address passed for padding: %s", addr))
}

func PadAddress(addr string) (string, error) {
	if len(addr) != 42 {
		return "", ErrInvalidAddress(addr)
	}
	padding := "0x000000000000000000000000"
	return padding + addr[2:], nil
}
