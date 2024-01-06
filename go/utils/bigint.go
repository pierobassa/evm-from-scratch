package utils

import "math/big"

// IsNegative returns true if the given big.Int is negative, false otherwise.
func IsNegative(x *big.Int) bool {
	return x.Bit(255) == 1
}

// OverflowingNeg returns the negation of the given big.Int.
// It handles overflow by wrapping around within 256 bits.
func OverflowingNeg(x *big.Int) *big.Int {
	negated := new(big.Int).Neg(x) // Negate x

	// Handle overflow (wrap around within 256 bits)
	negated.Mod(negated, Modulus)

	return negated
}
