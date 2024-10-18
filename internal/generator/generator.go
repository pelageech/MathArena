package generator

import (
	"math/rand/v2"

	"github.com/pelageech/matharena/internal/math"
)

type Difficulty int

const (
	Easy Difficulty = iota
	Medium
	Hard
)

func (d Difficulty) String() string {
	switch d {
	case Easy:
		return "Easy"
	case Medium:
		return "Medium"
	case Hard:
		return "Hard"
	}
	return "Unknown"
}

type Generator interface {
	Generate() math.ExpressionInt
	Difficulty() Difficulty
}

type EasyGenerator struct{}

func (g *EasyGenerator) Generate() math.ExpressionInt {
	n := 2 + rand.N(2)

	sum := math.Sum{}

	for range n {
		sum = append(sum, math.Num(rand.N(50)))
	}

	return sum
}

func (g *EasyGenerator) Difficulty() Difficulty {
	return Easy
}
