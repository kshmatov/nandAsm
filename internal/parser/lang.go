package parser

import "github.com/pkg/errors"

var (
	symbolTable = map[string]int{
		"SCREEN": screen,
		"KBD":    keyboard,
		"R0":     0,
		"R1":     1,
		"R2":     2,
		"R3":     3,
		"R4":     4,
		"R5":     5,
		"R6":     6,
		"R7":     7,
		"R8":     8,
		"R9":     9,
		"R10":    10,
		"R11":    11,
		"R12":    12,
		"R13":    13,
		"R14":    14,
		"R15":    15}

	ops = map[string]int{
		"0":   0b101010,
		"1":   0b111111,
		"-1":  0b111010,
		"D":   0b001100,
		"A":   0b110000,
		"M":   0b110000,
		"!D":  0b001101,
		"!A":  0b110001,
		"!M":  0b110001,
		"D+1": 0b011111,
		"A+1": 0b110111,
		"M+1": 0b110111,
		"D-1": 0b001110,
		"A-1": 0b110010,
		"M-1": 0b110010,
		"D+A": 0b000010,
		"A+D": 0b000010,
		"D-A": 0b010011,
		"A-D": 0b000111,
		"D&A": 0b000000,
		"D|A": 0b010101,
		"D+M": 0b000010,
		"M+D": 0b000010,
		"D-M": 0b010011,
		"M-D": 0b000111,
		"D&M": 0b000000,
		"D|M": 0b010101,
	}

	src = map[rune]int{
		'A': 0,
		'M': 1,
	}

	dst = map[rune]int{
		'M': 0b001,
		'D': 0b010,
		'A': 0b100,
	}

	jmp = map[string]int{
		"JGT": 0b001,
		"JEQ": 0b010,
		"JGE": 0b011,
		"JLT": 0b100,
		"JNE": 0b101,
		"JLE": 0b110,
		"JMP": 0b111,
	}
)

func getJump(s string) (uint16, error) {
	if v, ok := jmp[s]; ok {
		return uint16(v), nil
	}
	return 0, errors.Wrap(ErrUnknownOp, s)
}

func getOp(s string) (uint16, error) {
	if v, ok := ops[s]; ok {
		return uint16(v) << 6, nil
	}
	return 0, errors.Wrap(ErrUnknownOp, s)
}

func getSrc(s rune) (uint16, error) {
	if v, ok := src[s]; ok {
		return uint16(v) << 12, nil
	}
	return 0, errors.Wrap(ErrUnknownOp, string(s))
}

func getDst(s string) (uint16, error) {
	res := 0
	for _, c := range s {
		if v, ok := dst[c]; !ok {
			return 0, errors.Wrap(ErrUnknownOp, s)
		} else {
			if res&v != 0 {
				return 0, errors.Wrap(ErrUnknownOp, s)
			}
			res |= v
		}
	}
	return uint16(res) << 3, nil
}

func isJumpDest(s string) (string, bool) {
	l := len(s)
	if s[0] == '(' && s[l-1] == ')' {
		return s[1 : l-1], true
	}
	return "", false
}

func isMemLabel(s string) (string, bool) {
	if s[0] == '@' {
		return s, true
	}
	return "", false
}
