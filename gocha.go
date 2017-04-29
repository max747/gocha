package gocha

import (
	"math/rand"
	"regexp/syntax"
	"time"
)

type gocha struct {
	prog *syntax.Prog
}

type Gocha interface {
	Gen() string
}

func New(pattern string) (error, Gocha) {
	exp, err := syntax.Parse(pattern, syntax.Perl)
	if err != nil {
		return err, nil
	}

	prog, err := syntax.Compile(exp.Simplify())
	if err != nil {
		return err, nil
	}

	g := gocha{
		prog: prog,
	}

	return nil, g
}

func (g gocha) Gen() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	prog := g.prog
	pc := uint32(prog.Start)
	result := []rune{}

	rsInstRuneAny := []intRange{
		intRange{
			a: 0,
			b: 0x10ffff, // largest rune value
		},
	}

	rsInstRuneAnyNotNL := []intRange{
		intRange{
			a: 0,
			b: 9,
		},
		intRange{
			a: 11,
			b: 0x10ffff,
		},
	}

	inProgress := true

	for inProgress {

		switch prog.Inst[pc].Op {

		case syntax.InstMatch:
			inProgress = false

		case syntax.InstFail:
			inProgress = false

		case syntax.InstAlt:

			if i := rnd.Intn(2); i%2 == 1 {
				pc = prog.Inst[pc].Out
			} else {
				pc = prog.Inst[pc].Arg
			}

		case syntax.InstCapture:
			pc = prog.Inst[pc].Out

		case syntax.InstRuneAny:
			c := rune(randFromRange(rsInstRuneAny, rnd))
			result = append(result, c)
			pc = prog.Inst[pc].Out

		case syntax.InstRuneAnyNotNL:
			c := rune(randFromRange(rsInstRuneAnyNotNL, rnd))
			result = append(result, c)
			pc = prog.Inst[pc].Out

		case syntax.InstEmptyWidth:
			pc = prog.Inst[pc].Out

		case syntax.InstNop:
			return ""

		case syntax.InstRune1:
			result = append(result, prog.Inst[pc].Rune[0])
			pc = prog.Inst[pc].Out

		case syntax.InstRune:
			rs := make([]intRange, 0, len(prog.Inst[pc].Rune)/2)
			for i := 0; i < len(prog.Inst[pc].Rune); i += 2 {
				r := intRange{
					a: int(prog.Inst[pc].Rune[i]),
					b: int(prog.Inst[pc].Rune[i+1]),
				}
				rs = append(rs, r)
			}

			c := rune(randFromRange(rs, rnd))

			result = append(result, c)
			pc = prog.Inst[pc].Out

		default:
			panic("panic")
		}
	}

	return string(result)
}

type intRange struct {
	a int
	b int
}

func randFromRange(rs []intRange, rnd *rand.Rand) int {

	overallLen := 0

	for _, r := range rs {
		overallLen = overallLen + (r.b - r.a + 1)
	}
	index := rnd.Intn(overallLen)
	var result int
	for _, r := range rs {

		if (r.b - r.a) >= index {
			result = rnd.Intn(r.b-r.a+1) + r.a
			break
		}

		index = index - (r.b - r.a + 1)
	}

	return result
}
