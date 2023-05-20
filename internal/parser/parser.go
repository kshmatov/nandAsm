package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const (
	screen             = 0x4000
	keyboard           = 0x6000
	maxRamAddr         = keyboard
	maxRomAddr         = 0xFFFF
	firstUnassignedMem = 0x16

	lineFmt = "line <%v>: %v"
)

var (
	ErrEmptySource     = errors.New("no assembly code in source")
	ErrBadLine         = errors.New("bad format")
	ErrUnknownOp       = errors.New("unknown operator")
	ErrUnknownSymbol   = errors.New("unknown symbol")
	ErrLabelRedeclared = errors.New("label redeclared")
	ErrNotEnoughRam    = errors.New("not enough RAM")
	ErrNotEnoughRom    = errors.New("not enough ROM")

	firstFreeMem = firstUnassignedMem
	bindedMem    map[int]struct{}
)

func init() {
	bindedMem = make(map[int]struct{})
	for i := 0; i < 16; i++ {
		bindedMem[i] = struct{}{}
	}
}

type srcLine struct {
	n int
	l string
}

func (s srcLine) String() string {
	return fmt.Sprintf(lineFmt, s.n, s.l)
}

type MachineCode []uint16

func (m MachineCode) Binary() []byte {
	return nil
}

func (m MachineCode) String() []string {
	r := make([]string, len(m))
	for i, l := range m {
		r[i] = fmt.Sprintf("%016b", l)
	}
	return r
}

func Parse(source []string) (MachineCode, error) {
	if len(source) == 0 {
		return nil, ErrEmptySource
	}
	fSrc, err := firstPass(source)
	if err != nil {
		return nil, err
	}
	if len(fSrc) == 0 {
		return nil, ErrEmptySource
	}
	return secondPass(fSrc)
}

func firstPass(source []string) ([]srcLine, error) {
	prepared := make([]srcLine, 0, len(source))
	cnt := 1
	for i, line := range source {
		l := strings.TrimSpace(line)
		if l == "" {
			continue
		}
		parts := strings.Split(l, "//")
		l = parts[0]
		if l == "" {
			continue
		}
		if dst, ok := isJumpDest(l); ok {
			err := addJmpDest(dst, cnt+1)
			if err != nil {
				return nil, errors.Wrapf(err, lineFmt, i, line)
			}
			continue
		}
		cnt++
		if dst, ok := isMemLabel(l); ok {
			err := addMemLabel(dst)
			if err != nil {
				return nil, errors.Wrapf(err, lineFmt, i, line)
			}
		}
		prepared = append(prepared, srcLine{n: i, l: l})
	}
	return prepared, nil
}

func formatAOp(i int) uint16 {
	return 0b0111111111111111 & uint16(i)
}

func secondPass(sourse []srcLine) (MachineCode, error) {
	mc := make(MachineCode, len(sourse))
	var err error

	for lNum, line := range sourse {
		if line.l[0] == '@' {
			i, err := getMemLabel(line.l)
			if err != nil {
				return nil, errors.Wrap(ErrUnknownSymbol, line.String())
			}
			mc[lNum] = formatAOp(i)
			continue
		}
		p := strings.Split(line.l, ";")
		if len(p) > 2 {
			return nil, errors.Wrap(ErrBadLine, line.String())
		}
		var src, op, dst, jmp uint16
		if len(p) == 1 {
			jmp = 0
		} else {
			jmp, err = getJump(p[1])
			if err != nil {
				return nil, errors.Wrapf(err, line.String())
			}
		}
		p = strings.Split(p[0], "=")
		if len(p) > 2 {
			return nil, errors.Wrap(ErrBadLine, line.String())
		}
		var opStr string
		if len(p) == 2 {
			opStr = p[1]
			dst, err = getDst(p[0])
			if err != nil {
				return nil, errors.Wrap(err, line.String())
			}
		} else {
			opStr = p[0]
		}
		op, src, err = extractOpAndSrc(opStr)
		if err != nil {
			return nil, errors.Wrap(err, line.String())
		}
		mc[lNum] = 0b1110000000000000 | src | op | dst | jmp
	}

	return mc, nil
}

func extractOpAndSrc(s string) (uint16, uint16, error) {
	r := make([]rune, 0, len(s))
	var op, sr uint16
	var err error
	for _, c := range s {
		if c == ' ' || c == '\t' {
			continue
		}
		sr, _ = getSrc(c)
		r = append(r, c)
	}
	op, err = getOp(string(r))
	if err != nil {
		return 0, 0, ErrUnknownOp
	}
	return op, sr, nil
}

func addJmpDest(s string, i int) error {
	if i > maxRomAddr {
		return ErrNotEnoughRom
	}
	_, ok := symbolTable[s]
	if !ok {
		symbolTable[s] = i
		return nil
	}
	return errors.Wrap(ErrLabelRedeclared, s)
}

func addMemLabel(s string) error {
	_, err := strconv.Atoi(s[1:])
	if err == nil {
		return nil
	}

	// for isBinded(firstFreeMem) {
	firstFreeMem++
	if firstFreeMem > maxRamAddr {
		return ErrNotEnoughRam
	}
	// }

	if _, ok := symbolTable[s]; !ok {
		symbolTable[s] = firstFreeMem
	}
	return nil
}

func getMemLabel(s string) (int, error) {
	i, err := strconv.Atoi(s[1:])
	if err == nil {
		return i, nil
	}

	i, ok := symbolTable[s]
	if !ok {
		return 0, errors.Wrap(ErrUnknownSymbol, s)
	}
	return i, nil
}

func isBinded(i int) bool {
	_, ok := bindedMem[i]
	return ok
}

func bindMem(i int) {
	bindedMem[i] = struct{}{}
}
