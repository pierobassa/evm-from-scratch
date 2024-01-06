package opcodes

import (
	"evm-from-scratch-go/utils"
	"math/big"
)

// ApplyComparisonOp applies various comparison operations based on the provided opcode.
// Supported operations include lt, gt, slt, and sgt.
func ApplyComparisonOp(opcode OpCode, pc *int, stack *[]*big.Int) bool {
	if len(*stack) < 2 {
		return false
	}

	// Pop the last two elements from the stack.
	a, _ := popLastElement(pc, stack)
	b, _ := popLastElement(pc, stack)

	var result *big.Int
	switch opcode {
	case Lt:
		result = ltComparison(a, b)
	case Gt:
		result = gtComparison(a, b)
	case Slt:
		result = sltComparison(a, b)
	case Sgt:
		result = sgtComparison(a, b)
	default:
		return false
	}

	*stack = append(*stack, result)
	*pc++

	return true
}

// ltComparison returns 1 if a < b, 0 otherwise.
func ltComparison(a, b *big.Int) *big.Int {
	if a.Cmp(b) < 0 {
		return big.NewInt(1)
	}
	return big.NewInt(0)
}

// gtComparison returns 1 if a > b, 0 otherwise.
func gtComparison(a, b *big.Int) *big.Int {
	if a.Cmp(b) > 0 {
		return big.NewInt(1)
	}
	return big.NewInt(0)
}

// sltComparison returns 1 if a < b, 0 otherwise.
// It compares the two values as signed integers.
// If both a and b are negative, it compares their absolute values.
func sltComparison(a, b *big.Int) *big.Int {
	switch {
	case utils.IsNegative(a) && !utils.IsNegative(b):
		return big.NewInt(1)
	case !utils.IsNegative(a) && utils.IsNegative(b):
		return big.NewInt(0)
	case utils.IsNegative(a) && utils.IsNegative(b):
		aNeg := utils.OverflowingNeg(a)
		bNeg := utils.OverflowingNeg(b)
		if aNeg.Cmp(bNeg) <= 0 {
			return big.NewInt(0)
		}
		return big.NewInt(1)
	default:
		if a.Cmp(b) < 0 {
			return big.NewInt(1)
		}
		return big.NewInt(0)
	}
}

// sgtComparison returns 1 if a > b, 0 otherwise.
// It compares the two values as signed integers.
// If both a and b are negative, it compares their absolute values.
func sgtComparison(a, b *big.Int) *big.Int {
	switch {
	case utils.IsNegative(a) && !utils.IsNegative(b):
		return big.NewInt(0)
	case !utils.IsNegative(a) && utils.IsNegative(b):
		return big.NewInt(1)
	case utils.IsNegative(a) && utils.IsNegative(b):
		aNeg := utils.OverflowingNeg(a)
		bNeg := utils.OverflowingNeg(b)
		if aNeg.Cmp(bNeg) >= 0 {
			return big.NewInt(0)
		}
		return big.NewInt(1)
	default:
		if a.Cmp(b) > 0 {
			return big.NewInt(1)
		}
		return big.NewInt(0)
	}
}
