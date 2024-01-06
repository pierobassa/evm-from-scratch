package opcodes

import (
	"evm-from-scratch-go/pool"
	"evm-from-scratch-go/utils"
	"math/big"
)

// mod performs modulus operation with a check for division by zero.
// It calculates a % b and handles the edge case where b is zero.
// In EVM, MOD is used for modulo operation with unsigned integers.
func mod(a, b *big.Int) *big.Int {
	if b.Sign() == 0 {
		return big.NewInt(0) // Return 0 to avoid division by zero
	}
	result := pool.FetchBigIntFromPool()
	defer pool.ReleaseBigIntToPool(result)
	return result.Mod(a, b)
}

// div performs division, with a check for division by zero.
// It calculates a / b and returns 0 when b is zero, adhering to EVM rules.
// In EVM, DIV is used for division of unsigned integers.
func div(a, b *big.Int) *big.Int {
	if b.Sign() == 0 {
		return big.NewInt(0) // Return 0 to avoid division by zero
	}
	result := pool.FetchBigIntFromPool()
	defer pool.ReleaseBigIntToPool(result)
	return result.Div(a, b)
}

// ApplyArithmeticOp applies various arithmetic operations based on the provided opcode.
// Supported operations include add, sub, mul, div, mod, addmod, mulmod, and exp.
// It pops the required operands from the stack, performs the operation, and pushes the result back.
func ApplyArithmeticOp(opcode OpCode, pc *int, stack *[]*big.Int) bool {
	if len(*stack) < 2 {
		return false
	}

	// Pop the last two elements from the stack.
	a, _ := popLastElement(pc, stack)
	b, _ := popLastElement(pc, stack)

	result := pool.FetchBigIntFromPool()
	defer pool.ReleaseBigIntToPool(result)

	switch opcode {
	case Add:
		result.Add(a, b)
	case Sub:
		result.Sub(a, b)
	case Mul:
		result.Mul(a, b)
	case Div:
		result = div(a, b)
	case Mod:
		result = mod(a, b)
	case Addmod:
		intermediate := pool.FetchBigIntFromPool().Add(a, b)
		defer pool.ReleaseBigIntToPool(intermediate) // Release intermediate to the pool after use. defer ensures this is called even if the function panics.

		modValue, _ := popLastElement(pc, stack)

		result.Mod(intermediate, modValue)
	case Mulmod:
		intermediate := pool.FetchBigIntFromPool().Mul(a, b)
		defer pool.ReleaseBigIntToPool(intermediate)

		modValue, _ := popLastElement(pc, stack)
		result.Mod(intermediate, modValue)
	case Exp:
		result.Exp(a, b, nil)
	default:
		return false
	}

	// Apply modulus to keep result within 256 bits
	result.Mod(result, utils.Modulus)

	*stack = append(*stack, result)
	*pc++

	return true
}

// SignedDivision performs signed division of the two top elements of the stack.
// It handles division by zero and sign considerations.
func SignedDivision(pc *int, stack *[]*big.Int) bool {
	if len(*stack) < 2 {
		return false
	}

	a, _ := popLastElement(pc, stack)
	b, _ := popLastElement(pc, stack)

	if b.Sign() == 0 {
		*stack = append(*stack, big.NewInt(0))
		*pc++
		return true
	}

	result := applySignedDivision(a, b)
	*stack = append(*stack, result)
	*pc++
	return true
}

// applySignedDivision performs the actual signed division logic.
// It is extracted for clarity and reusability.
func applySignedDivision(a, b *big.Int) *big.Int {
	result := pool.FetchBigIntFromPool()
	defer pool.ReleaseBigIntToPool(result)

	aIsNegative := utils.IsNegative(a)
	bIsNegative := utils.IsNegative(b)

	if aIsNegative {
		a = utils.OverflowingNeg(a)
	}
	if bIsNegative {
		b = utils.OverflowingNeg(b)
	}

	result.Div(a, b)

	// Apply sign
	if aIsNegative != bIsNegative {
		result.Neg(result)
	}

	result.Mod(result, utils.Modulus)
	return result
}

// SignedModulus performs signed modulus of the two top elements of the stack.
// Similar to SignedDivision, it handles special cases for division by zero and sign.
func SignedModulus(pc *int, stack *[]*big.Int) bool {
	if len(*stack) < 2 {
		return false
	}

	a, _ := popLastElement(pc, stack)
	b, _ := popLastElement(pc, stack)

	if b.Sign() == 0 {
		*stack = append(*stack, big.NewInt(0))
		*pc++
		return true
	}

	result := applySignedModulus(a, b)
	*stack = append(*stack, result)
	*pc++
	return true
}

// applySignedModulus performs the actual signed modulus logic.
func applySignedModulus(a, b *big.Int) *big.Int {
	result := pool.FetchBigIntFromPool()
	defer pool.ReleaseBigIntToPool(result)

	aIsNegative := utils.IsNegative(a)
	bIsNegative := utils.IsNegative(b)

	if aIsNegative {
		a = utils.OverflowingNeg(a)
	}
	if bIsNegative {
		b = utils.OverflowingNeg(b)
	}

	result.Mod(a, b)

	// Apply sign
	if aIsNegative {
		result.Neg(result)
	}

	result.Mod(result, utils.Modulus)
	return result
}

// SignedExtend extends the size of a byte within a 256-bit word based on the value of 'k'.
// This operation is aligned with the EVM's SIGNEXTEND opcode as described in the Ethereum Yellow Paper.
// The function pops two elements from the stack: 'k' and 'x', then pushes the result after sign extension.
func SignedExtend(pc *int, stack *[]*big.Int) bool {
	if len(*stack) < 2 {
		return false
	}

	k, _ := popLastElement(pc, stack)
	x, _ := popLastElement(pc, stack)

	if k.Cmp(big.NewInt(31)) > 0 {
		*stack = append(*stack, x)
		*pc++
		return true
	}

	result := applySignedExtend(x, k)
	*stack = append(*stack, result)
	*pc++

	return true
}

// applySignedExtend handles the core logic of the SIGNEXTEND operation.
// It extends the sign of the byte at the position 'k' to the left within the 256-bit word 'x'.
// This operation is crucial for correctly interpreting signed integers in EVM.
func applySignedExtend(x, k *big.Int) *big.Int {
	result := pool.FetchBigIntFromPool()
	defer pool.ReleaseBigIntToPool(result)

	bytes := x.Bytes() // Convert 'x' to a byte slice

	// Padding 'x' to 32 bytes (256 bits) if necessary
	for len(bytes) < 32 {
		bytes = append(bytes, 0)
	}

	byteIndex := int(k.Uint64()) // Determine the index of the byte to extend
	signByte := bytes[byteIndex] // Fetch the sign byte

	// Extend the sign of the byte at 'byteIndex' to the left of the word
	for i := 0; i < 32; i++ {
		if i > byteIndex {
			/*
				Checking with 0x7f:

				0x7f in binary is 01111111. This is the largest signed 8-bit number.
				When you compare the sign byte with 0x7f, you are effectively checking if the sign bit is set (i.e., if the byte represents a negative number in two's complement).
				If the sign byte is greater than 0x7f, its most significant bit must be 1, indicating it's negative
			*/
			if signByte > 0x7f {
				/*
					If the number is negative (sign bit is 1), you extend the sign by setting all more significant bytes to 0xFF (11111111 in binary).
					This preserves the negative sign across the entire 256-bit word.
				*/
				bytes[i] = 0xFF
			} else {
				/*
					Conversely, if the number is positive (or zero), you extend the sign by setting all more significant bytes to 0x00. This effectively pads the number with zeros on the left, preserving its positive value.
				*/
				bytes[i] = 0x00
			}
		}
	}

	// Convert the byte slice back to big.Int and return
	result.SetBytes(utils.FromLittleEndian(bytes))

	return result
}
