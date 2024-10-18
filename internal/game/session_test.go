package game

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/pelageech/matharena/internal/game/generator/mocks"
	"github.com/pelageech/matharena/internal/game/math"
)

type clock struct {
	timeNow time.Time
}

func (c *clock) now() time.Time {
	return c.timeNow
}

func (c *clock) add(t time.Duration) {
	c.timeNow = c.timeNow.Add(t)
}

func TestSession(t *testing.T) {
	clck := &clock{}
	generator := mocks.NewGenerator(t)
	generator.EXPECT().Generate().Return(math.Num(10)).Once()
	generator.EXPECT().Generate().Return(math.Num(20)).Once()
	generator.EXPECT().Generate().Return(math.Num(30)).Once()
	generator.EXPECT().Generate().Return(math.Num(40)).Once()

	s := NewSession(42, time.Second, generator, clck.now(), WithDeltas(Deltas{
		OnCorrect:   100 * time.Millisecond,
		OnIncorrect: 100 * time.Millisecond,
	}))

	t.Run("10,correct", func(t *testing.T) {
		clck.add(200 * time.Millisecond)
		s.Answer(10, clck.now())
		assert.Equal(t, 900*time.Millisecond, s.timeLeft)
	})

	t.Run("20,correct", func(t *testing.T) {
		clck.add(500 * time.Millisecond)
		s.Answer(20, clck.now())
		assert.Equal(t, 500*time.Millisecond, s.timeLeft)
	})

	t.Run("30,incorrect", func(t *testing.T) {
		clck.add(100 * time.Millisecond)
		s.Answer(42, clck.now())
		assert.Equal(t, 300*time.Millisecond, s.timeLeft)
	})

	t.Run("40,time left", func(t *testing.T) {
		clck.add(300 * time.Millisecond)
		s.Answer(40, clck.now())
		assert.Equal(t, time.Duration(0), s.timeLeft)
		assert.Equal(t, clck.now(), s.finishTime)
	})
}
