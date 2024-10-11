package generator

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"strconv"

	"github.com/pelageech/matharena/internal/math"
)

var ErrDifficultyNotRecognized = errors.New("difficulty not recognized")

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

type Generator struct {
	difficulty Difficulty
}

func New(difficulty Difficulty) *Generator {
	return &Generator{difficulty}
}

func (g *Generator) Generate() (math.ExpressionInt, error) {
	switch g.difficulty {
	case Easy:
		return g.generateEasy(), nil
	}
	return nil, fmt.Errorf("%s: %w", strconv.Quote(g.difficulty.String()), ErrDifficultyNotRecognized)
}

func (g *Generator) generateEasy() math.ExpressionInt {
	n := 2 + rand.N(2)

	sum := math.Sum{}

	for range n {
		sum = append(sum, math.Num(rand.N(50)))
	}

	return sum
}
