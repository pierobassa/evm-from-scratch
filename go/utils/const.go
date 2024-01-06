package utils

import "math/big"

// Modulus to wrap big integers to 256 bits.
var Modulus = new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)
