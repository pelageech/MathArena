package game

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"strconv"
	"time"

	"github.com/pelageech/matharena/internal/game/generator"
	"github.com/pelageech/matharena/internal/game/math"
)

const (
	_defaultDeltaOnCorrect   = 5 * time.Second
	_defaultDeltaOnIncorrect = 5 * time.Second
)

type SessionID int64

func (id SessionID) String() string {
	return strconv.FormatInt(int64(id), 10)
}

func ParseSessionID(s string) (SessionID, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid session id: %s", s)
	}
	return SessionID(i), nil
}

func newSessionID() SessionID {
	return SessionID(rand.Int64())
}

type Deltas struct {
	OnCorrect   time.Duration
	OnIncorrect time.Duration
}

type Session struct {
	sessionID         SessionID
	userID            int
	currentExpression math.ExpressionInt
	answer            int
	score             int
	timeLeft          time.Duration
	generator         generator.Generator

	startTime            time.Time
	finishTime           time.Time
	lastUpdateExpression time.Time
	deltas               Deltas
}

type Opt func(*Session)

func WithDeltas(deltas Deltas) Opt {
	return func(s *Session) {
		s.deltas = deltas
	}
}

func WithCustomID(id SessionID) Opt {
	return func(s *Session) {
		s.sessionID = id
	}
}

func NewSession(userID int, timeStart time.Duration, generator generator.Generator, timeNow time.Time, opts ...Opt) *Session {
	s := &Session{
		sessionID: newSessionID(),
		userID:    userID,
		timeLeft:  timeStart,
		generator: generator,
		startTime: timeNow,
		deltas: Deltas{
			OnCorrect:   _defaultDeltaOnCorrect,
			OnIncorrect: _defaultDeltaOnIncorrect,
		},
	}

	for _, opt := range opts {
		opt(s)
	}

	s.updateExpression(timeNow)
	return s
}

var (
	ErrAnswerIsIncorrect = errors.New("answer is incorrect")
	ErrTimeIsLeft        = errors.New("time is left")
)

// Answer handles answer from the user. Checks if the session is ended before it.
func (s *Session) Answer(answer int, timeNow time.Time) error {
	s.updateTimeOnAnswer(timeNow)
	// check if the user is late to answer
	if s.timeLeft <= 0 {
		s.Stop(timeNow.Add(s.timeLeft))
		return ErrTimeIsLeft
	}

	defer s.updateExpression(timeNow)

	if answer != s.answer {
		s.timeOnIncorrect()
		if s.timeLeft <= 0 { // check after incorrect answer
			s.Stop(timeNow.Add(s.timeLeft))
			return ErrTimeIsLeft
		}
		return ErrAnswerIsIncorrect
	}

	s.timeOnCorrect()
	s.UpdateScore()

	return nil
}

func (s *Session) CheckTime(timeNow time.Time) bool {
	if timeNow.Before(s.lastUpdateExpression.Add(s.timeLeft)) {
		return true
	}

	s.Stop(s.lastUpdateExpression.Add(s.timeLeft))
	return false
}

func (s *Session) Stop(finishTime time.Time) {
	s.finishTime = finishTime
}

func (s *Session) UpdateScore() {
	s.score++
}

func (s *Session) ID() SessionID {
	return s.sessionID
}

func (s *Session) TimeLeft() time.Duration {
	return s.timeLeft
}

func (s *Session) FinishTime() time.Time {
	return s.finishTime
}

func (s *Session) CurrentExpression() math.ExpressionInt {
	return s.currentExpression
}

func (s *Session) UserID() int {
	return s.userID
}

func (s *Session) Score() int {
	return s.score
}

func (s *Session) timeOnCorrect() {
	s.timeLeft += s.deltas.OnCorrect
}

func (s *Session) timeOnIncorrect() {
	s.timeLeft -= s.deltas.OnIncorrect
}

func (s *Session) updateTimeOnAnswer(timeNow time.Time) {
	s.timeLeft -= timeNow.Sub(s.lastUpdateExpression)
}

func (s *Session) updateExpression(timeNow time.Time) {
	s.currentExpression = s.generator.Generate()
	s.answer = s.currentExpression.Calculate()
	s.lastUpdateExpression = timeNow
}
