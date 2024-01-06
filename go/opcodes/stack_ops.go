package opcodes

import "math/big"

// PushX reads 'size' bytes from the EVM code starting from the current program counter (PC)
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
func PushX(pc *int, stack *[]*big.Int, code []byte, size int) {
	end := *pc + size
	if end <= len(code) {
		value := new(big.Int).SetBytes(code[*pc:end])

		// Prepend the value to the stack.
		*stack = append(*stack, value)

		*pc = end
	}
}

// PopX pops elements from the stack.
// It takes the program counter, stack, and number of elements to pop as input.
// It returns the popped elements and a success indicator.
func PopX(pc *int, stack *[]*big.Int, size int) ([]*big.Int, bool) {
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

// popLastElement pops the last element from the stack.
// It takes the program counter and stack as input.
// It returns the popped element and a success indicator.
func popLastElement(pc *int, stack *[]*big.Int) (*big.Int, bool) {
	if len(*stack) < 1 {
		return nil, false
	}

	elements, _ := PopX(pc, stack, 1)

	return elements[0], true
}
