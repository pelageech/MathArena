package game

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/pelageech/matharena/internal/game/generator"
	"github.com/pelageech/matharena/internal/game/math"
)

const (
	_defaultDeltaOnCorrect   = 5 * time.Second
	_defaultDeltaOnIncorrect = 5 * time.Second
)

type uuidString = string

type Deltas struct {
	OnCorrect   time.Duration
	OnIncorrect time.Duration
}

type Session struct {
	sessionID         uuidString
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

func NewSession(userID int, timeStart time.Duration, generator generator.Generator, timeNow time.Time, opts ...Opt) *Session {
	s := &Session{
		sessionID: uuid.NewString(),
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

	s.updateExpression()
	return s
}

var (
	ErrAnswerIsIncorrect = errors.New("answer is incorrect")
	ErrTimeIsLeft        = errors.New("time is left")
)

func (s *Session) Answer(answer int, timeNow time.Time) error {
	s.updateTimeOnAnswer(timeNow)
	// check if the user is late to answer
	if !s.CheckTime(timeNow) {
		return ErrTimeIsLeft
	}

	defer s.updateExpression()

	if answer != s.answer {
		s.timeOnIncorrect()
		if !s.CheckTime(timeNow) { // check after incorrect answer
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

func (s *Session) timeOnCorrect() {
	s.timeLeft += s.deltas.OnCorrect
}

func (s *Session) timeOnIncorrect() {
	s.timeLeft -= s.deltas.OnIncorrect
}

func (s *Session) updateTimeOnAnswer(timeNow time.Time) {
	s.timeLeft -= timeNow.Sub(s.lastUpdateExpression)
	s.lastUpdateExpression = timeNow
}

func (s *Session) updateExpression() {
	s.currentExpression = s.generator.Generate()
	s.answer = s.currentExpression.Calculate()
}
