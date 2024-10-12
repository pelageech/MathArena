package math

import (
	"bytes"
	"strconv"
)

type ExpressionInt interface {
	Calculate() int
	Marshal() []byte
}

type Num int

func (n Num) Calculate() int {
	return int(n)
}

func (n Num) Marshal() []byte {
	return []byte(strconv.Itoa(int(n)))
}

type Fact struct {
	Fact ExpressionInt `json:"fact"`
}

func (n Fact) Calculate() int {
	f := 1
	for i := 1; i <= n.Fact.Calculate(); i++ {
		f *= i
	}
	return f
}

func (n Fact) Marshal() []byte {
	return append(n.Fact.Marshal(), '!')
}

type Sum []ExpressionInt

func (s Sum) Calculate() int {
	i := 0
	for _, x := range s {
		i += x.Calculate()
	}
	return i
}

func (s Sum) Marshal() []byte {
	if len(s) == 0 {
		return nil
	}

	b := &bytes.Buffer{}

	switch v := s[0].(type) {
	case Sum, Num:
		b.Write(v.Marshal())
	default:
		b.WriteByte('(')
		b.Write(v.Marshal())
		b.WriteByte(')')
	}

	for _, x := range s[1:] {
		switch v := x.(type) {
		case Sum:
			b.WriteByte('+')
			b.Write(v.Marshal())
		case Num:
			if v.Calculate() < 0 {
				b.Write(v.Marshal())
			} else {
				b.WriteByte('+')
				b.Write(v.Marshal())
			}
		default:
			b.WriteByte('+')
			b.WriteByte('(')
			b.Write(v.Marshal())
			b.WriteByte(')')
		}
	}

	return b.Bytes()
}
