package evm

import (
	"evm-from-scratch-go/opcodes"
	"math/big"
	"slices"
)

// Evm executes the EVM code and returns the stack and a success indicator.
// It takes the EVM code as input and returns the stack and a success indicator.
// The stack is returned in reverse order, with the top element at the end.
// The success indicator is true if the execution was successful, false otherwise.
func Evm(code []byte) ([]*big.Int, bool) {
	var stack []*big.Int
	pc := 0

	for pc < len(code) {
		opcode, err := opcodes.NewOpCode(code[pc])

		if err != nil {
			return nil, false // Revert on unknown opcodes
		}

		pc++

		// Stop execution if the opcodes is STOP
		if opcode == opcodes.Stop {
			return stack, true // Halt execution
		}

		if !executeOpcode(&pc, &stack, code, opcode) {
			return nil, false
		}
	}

	// Reverse the stack so that the top element is at the end.
	slices.Reverse(stack)

	return stack, true // Success
}

// executeOpcode executes the opcode and returns true if the execution was successful, false otherwise.
// It takes the program counter, stack, EVM code, and opcode as input.
// It executes the opcode and returns true if the execution was successful, false otherwise.
func executeOpcode(pc *int, stack *[]*big.Int, code []byte, opcode opcodes.OpCode) bool {
	switch opcode {
	case opcodes.Push0:
		*stack = append(*stack, big.NewInt(0)) // Push 0 onto the stack
	case opcodes.Push1, opcodes.Push2, opcodes.Push4, opcodes.Push6, opcodes.Push10, opcodes.Push11, opcodes.Push32:
		opcodes.PushX(pc, stack, code, opcodes.PushOpcodeToBytes[opcode])
	case opcodes.Pop:
		if _, ok := opcodes.PopX(pc, stack, 1); !ok {
			return false
		}
	case opcodes.Add, opcodes.Sub, opcodes.Mul, opcodes.Div, opcodes.Mod, opcodes.Addmod, opcodes.Mulmod, opcodes.Exp:
		if !opcodes.ApplyArithmeticOp(opcode, pc, stack) {
			return false
		}
	case opcodes.Signextend:
		if !opcodes.SignedExtend(pc, stack) {
			return false
		}
	case opcodes.Sdiv:
		if !opcodes.SignedDivision(pc, stack) {
			return false
		}
	case opcodes.Smod:
		if !opcodes.SignedModulus(pc, stack) {
			return false
		}
	case opcodes.Lt, opcodes.Gt, opcodes.Slt, opcodes.Sgt:
		if !opcodes.ApplyComparisonOp(opcode, pc, stack) {
			return false
		}
	default:
		return false // Revert on unknown opcodes
	}

	return true
}
