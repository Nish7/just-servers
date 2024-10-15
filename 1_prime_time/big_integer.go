package main

import (
	"fmt"
	"math/big"
)

type BigNumber struct {
	BigInt  *big.Int
	IsFloat bool
}

func (n *BigNumber) UnmarshalJSON(data []byte) error {
	// Try to unmarshal into a string first
	numStr := string(data)

	floatValue := new(big.Float)
	if _, ok := floatValue.SetString(numStr); ok {
		if floatValue.IsInt() {
			n.BigInt = new(big.Int)
			floatValue.Int(n.BigInt)
			n.IsFloat = false
		} else {
			n.IsFloat = true
		}
	} else {
		return fmt.Errorf("invalid number format: %s", numStr)
	}
	return nil
}
