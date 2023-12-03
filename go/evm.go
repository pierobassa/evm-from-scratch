package evm

import (
	"math/big"
	"slices"
)

// Modulus to wrap big integers to 256 bits.
var modulus = new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)

// Evm executes the EVM code and returns the stack and a success indicator.
func Evm(code []byte) ([]*big.Int, bool) {
	var stack []*big.Int
	pc := 0

	for pc < len(code) {
		opcode, err := NewOpCode(code[pc])
		if err != nil {
			return nil, false // Revert on unknown opcode
		}
		pc++

		switch opcode {
		case Stop:
			return stack, true // Halt execution
		case Push0:
			stack = append(stack, big.NewInt(0)) // Push 0 onto the stack
		case Push1:
			pushX(&pc, &stack, code, 1)
		case Push2:
			pushX(&pc, &stack, code, 2)
		case Push4:
			pushX(&pc, &stack, code, 4)
		case Push6:
			pushX(&pc, &stack, code, 6)
		case Push10:
			pushX(&pc, &stack, code, 10)
		case Push11:
			pushX(&pc, &stack, code, 11)
		case Push32:
			pushX(&pc, &stack, code, 32)
		case Pop:
			if !popX(&pc, &stack, 1) {
				return nil, false
			}
		case Add, Sub, Mul, Div:
			if !applyArithmeticOp(opcode, &pc, &stack) {
				return nil, false
			}
		}
	}

	// Reverse the stack so that the top element is at the end.
	slices.Reverse(stack)

	return stack, true // Success
}

// pushX reads 'size' bytes from the EVM code starting from the current program counter (PC)
// and pushes them as a big integer onto the EVM stack.
//
// Arguments:
// pc    - Pointer to the program counter which indicates the current position in the code.
// stack - Pointer to the EVM stack where all computational values are stored.
// code  - Byte slice representing the EVM code being executed.
// size  - Number of bytes to read from the code and push onto the stack.
//
// If the bytes to be read exceed the bounds of the code slice, no action is performed.
// After successfully reading the bytes and pushing them onto the stack, the program counter
// is updated to the position after the read bytes.
func pushX(pc *int, stack *[]*big.Int, code []byte, size int) {
	end := *pc + size
	if end <= len(code) {
		value := new(big.Int).SetBytes(code[*pc:end])

		// Prepend the value to the stack.
		*stack = append(*stack, value)

		*pc = end
	}
}

// popX pops elements from the stack.
func popX(pc *int, stack *[]*big.Int, size int) bool {
	if len(*stack) < size {
		return false
	}

	// Remove the last 'size' elements from the stack.
	*stack = (*stack)[:len(*stack)-size]
	*pc++
	return true
}

// applyArithmeticOp applies arithmetic operations like add, sub,0 mul, div.
func applyArithmeticOp(opcode OpCode, pc *int, stack *[]*big.Int) bool {
	if len(*stack) < 2 {
		return false
	}

	// Pop the last two elements from the stack.
	a, b := (*stack)[len(*stack)-1], (*stack)[len(*stack)-2]
	*stack = (*stack)[:len(*stack)-2]

	var result *big.Int
	switch opcode {
	case Add:
		result = new(big.Int).Add(a, b)
	case Sub:
		result = new(big.Int).Sub(a, b)
	case Mul:
		result = new(big.Int).Mul(a, b)
	case Div:
		if b.Sign() == 0 { // Check for division by zero
			result = big.NewInt(0)
		} else {
			result = new(big.Int).Div(a, b)
		}
	default:
		return false
	}

	// Apply modulus to keep result within 256 bits
	result.Mod(result, modulus)

	*stack = append(*stack, result)
	*pc++

	return true
}
