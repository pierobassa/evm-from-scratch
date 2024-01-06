package opcodes

import (
	"errors"
)

// OpCode represents the operation codes (opcodes) in the EVM.
type OpCode byte

// Define opcodes as constants.
const (
	Stop OpCode = iota
	Push0
	Push1
	Push2
	Push4
	Push6
	Push10
	Push11
	Push32
	Pop
	Add
	Mul
	Sub
	Div
	Sdiv
	Mod
	Smod
	Addmod
	Mulmod
	Exp
	Signextend
	Lt
	Gt
	Slt
	Sgt
)

// opCodeMap maps byte values to OpCode.
var opCodeMap = map[byte]OpCode{
	0:   Stop,
	95:  Push0,
	96:  Push1,
	97:  Push2,
	99:  Push4,
	101: Push6,
	105: Push10,
	106: Push11,
	127: Push32,
	80:  Pop,
	1:   Add,
	2:   Mul,
	3:   Sub,
	4:   Div,
	5:   Sdiv,
	6:   Mod,
	7:   Smod,
	8:   Addmod,
	9:   Mulmod,
	10:  Exp,
	11:  Signextend,
	16:  Lt,
	17:  Gt,
	18:  Slt,
	19:  Sgt,
}

// PushOpcodeToBytes maps number of bytes to read from the code and push onto the stack.
var PushOpcodeToBytes = map[OpCode]int{
	Push1:  1,
	Push2:  2,
	Push4:  4,
	Push6:  6,
	Push10: 10,
	Push11: 11,
	Push32: 32,
}

// NewOpCode tries to convert a byte into an OpCode. It returns an error if the opcodes is unknown.
func NewOpCode(b byte) (OpCode, error) {
	opcode, ok := opCodeMap[b]
	if !ok {
		return 0, errors.New("unknown opcodes")
	}
	return opcode, nil
}
