package evm

import (
	"math/big"
	"slices"
)

// Modulus to wrap big integers to 256 bits.
var modulus = new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)

// Evm executes the EVM code and returns the stack and a success indicator.
// TODO: Reduce cognitive complexity
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
			if _, ok := popX(&pc, &stack, 1); !ok {
				return nil, false
			}
		case Add, Sub, Mul, Div, Mod, Addmod, Mulmod, Exp:
			if !applyArithmeticOp(opcode, &pc, &stack) {
				return nil, false
			}
		case Signextend:
			if !signExtend(&pc, &stack) {
				return nil, false
			}
		case Sdiv:
			if !sdiv(&pc, &stack) {
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
func popX(pc *int, stack *[]*big.Int, size int) ([]*big.Int, bool) {
	if len(*stack) < size {
		return nil, false
	}

	// Get the last 'size' elements from the stack.
	elements := (*stack)[len(*stack)-size:]

	// Remove the last 'size' elements from the stack.
	*stack = (*stack)[:len(*stack)-size]
	*pc++

	return elements, true
}

func mod(a, b *big.Int) *big.Int {
	if b.Sign() == 0 { // Check for division by zero
		return big.NewInt(0)
	}
	return new(big.Int).Mod(a, b)
}

func div(a, b *big.Int) *big.Int {
	if b.Sign() == 0 { // Check for division by zero
		return big.NewInt(0)
	}
	return new(big.Int).Div(a, b)
}

func popLastElement(pc *int, stack *[]*big.Int) (*big.Int, bool) {
	if len(*stack) < 1 {
		return nil, false
	}

	elements, _ := popX(pc, stack, 1)

	return elements[0], true
}

// applyArithmeticOp applies arithmetic operations like add, sub, mul, div, mod
func applyArithmeticOp(opcode OpCode, pc *int, stack *[]*big.Int) bool {
	if len(*stack) < 2 {
		return false
	}

	// Pop the last two elements from the stack.
	a, _ := popLastElement(pc, stack)
	b, _ := popLastElement(pc, stack)

	var result *big.Int
	switch opcode {
	case Add:
		result = new(big.Int).Add(a, b)
	case Sub:
		result = new(big.Int).Sub(a, b)
	case Mul:
		result = new(big.Int).Mul(a, b)
	case Div:
		result = div(a, b)
	case Mod:
		result = mod(a, b)
	case Addmod:
		result = new(big.Int).Add(a, b)
		modValue, _ := popLastElement(pc, stack) // Get the last element from the stack
		result = mod(result, modValue)
	case Mulmod:
		result = new(big.Int).Mul(a, b)
		modValue, _ := popLastElement(pc, stack) // Get the last element from the stack
		result = mod(result, modValue)
	case Exp:
		result = new(big.Int).Exp(a, b, nil)
	default:
		return false
	}

	// Apply modulus to keep result within 256 bits
	result.Mod(result, modulus)

	*stack = append(*stack, result)
	*pc++

	return true
}

// fromLittleEndian converts a byte slice from little endian to big endian.
func fromLittleEndian(bytes []byte) []byte {
	copyBytes := bytes

	for i := 0; i < len(copyBytes)/2; i++ {
		copyBytes[i], copyBytes[len(copyBytes)-i-1] = copyBytes[len(copyBytes)-i-1], copyBytes[i]
	}

	return copyBytes
}

func signExtend(pc *int, stack *[]*big.Int) bool {
	if len(*stack) < 2 {
		return false
	}

	// Pop the last two elements from the stack.
	k, _ := popLastElement(pc, stack)
	x, _ := popLastElement(pc, stack)

	// If k is greater than 31, push x back onto the stack
	if k.Cmp(big.NewInt(31)) > 0 {
		*stack = append(*stack, x)
		*pc++

		return true
	}

	// Convert x to a byte slice
	bytes := x.Bytes()

	// Extend to 32 bytes if necessary
	for len(bytes) < 32 {
		bytes = append(bytes, 0)
	}

	// Perform the sign extension
	byteIndex := int(k.Uint64())
	signByte := bytes[byteIndex]

	for i := 0; i < 32; i++ {
		if i > int(k.Uint64()) {
			if signByte > 0x7f { // Check if the sign bit is set
				bytes[i] = 0xFF
			} else {
				bytes[i] = 0x00
			}
		}

	}

	// Convert the byte slice back to big.Int
	result := new(big.Int).SetBytes(fromLittleEndian(bytes))
	*stack = append(*stack, result)

	*pc++

	return true
}

// Helper function to check if the most significant bit is set (negative in two's complement)
func isNegative(x *big.Int) bool {
	return x.Bit(255) == 1
}

// Helper function to negate a big.Int and handle overflow
func overflowingNeg(x *big.Int) *big.Int {
	negated := new(big.Int).Neg(x) // Negate x

	// Handle overflow (wrap around within 256 bits)
	negated.Mod(negated, modulus)

	return negated
}

func sdiv(pc *int, stack *[]*big.Int) bool {
	if len(*stack) < 2 {
		return false
	}

	a, _ := popLastElement(pc, stack)
	b, _ := popLastElement(pc, stack)

	if b.Sign() == 0 { // Division by zero
		*stack = append(*stack, big.NewInt(0))
		*pc++
		return true
	}

	aIsNegative := isNegative(a)
	bIsNegative := isNegative(b)

	if aIsNegative {
		a = overflowingNeg(a)
	}
	if bIsNegative {
		b = overflowingNeg(b)
	}

	result := new(big.Int).Div(a, b)

	// Apply sign
	if aIsNegative != bIsNegative {
		result.Neg(result)
	}

	// Apply modulus to keep result within 256 bits
	result.Mod(result, modulus)

	*stack = append(*stack, result)
	*pc++

	return true
}
