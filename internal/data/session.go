package data

import (
	"context"
	"errors"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/pelageech/matharena/internal/game"
	"github.com/pelageech/matharena/internal/game/generator"
	"time"
)

type GameSessionsDB interface {
	CreateSession(ctx context.Context, userId int, startTime time.Time) (game.SessionID, error)
	FinishSession(ctx context.Context, userID int, id game.SessionID, finishTime time.Time) error
}

type SessionDataLayer struct {
	logger         *log.Logger
	db             GameSessionsDB
	activeSessions *game.ActiveSessionsPool
}

func NewSessionDataLayer(db GameSessionsDB, logger *log.Logger) *SessionDataLayer {
	return &SessionDataLayer{
		db:             db,
		activeSessions: game.NewActiveSessionsPool(),
		logger:         logger,
	}
}

func (ld *SessionDataLayer) CreateSession(ctx context.Context, timeStart time.Duration, userID int, _ generator.Difficulty, timeNow time.Time) (*game.Session, error) {
	id, err := ld.db.CreateSession(ctx, userID, timeNow)
	if err != nil {
		return nil, fmt.Errorf("db create session: %w", err)
	}

	s := game.NewSession(userID, timeStart, &generator.EasyGenerator{}, timeNow, game.WithCustomID(id))
	if err := ld.activeSessions.Put(s); err != nil {
		return nil, fmt.Errorf("set session active: %w", err)
	}
	return s, nil
}

func (ld *SessionDataLayer) Answer(ctx context.Context, sessionID game.SessionID, answer, userID int, timeNow time.Time) (*game.Session, error) {
	s, err := ld.activeSessions.Get(sessionID, timeNow)
	if err != nil {
		return nil, fmt.Errorf("get session %v: %w", sessionID, err)
	}

	if s.UserID() != userID {
		ld.logger.Infof("%v %v", s.UserID(), userID)
		return nil, errors.New("invalid session")
	}

	err = s.Answer(answer, timeNow)
	if errors.Is(err, game.ErrAnswerIsIncorrect) {
		ld.logger.Debugf("sid: %v, ans: %v: answer is incorrect", sessionID, answer)
	} else if err != nil {
		_ = ld.Stop(ctx, sessionID, userID, timeNow)
		return nil, fmt.Errorf("answer: %w", err)
	}

	return s, nil
}

func (ld *SessionDataLayer) Stop(ctx context.Context, sessionID game.SessionID, userID int, timeNow time.Time) error {
	s, err := ld.activeSessions.Get(sessionID, timeNow)
	if errors.Is(err, game.ErrTimeIsLeft) {
		return fmt.Errorf("sid %v: already stopped", sessionID)
	}
	if err != nil {
		return fmt.Errorf("get session %v: %w", sessionID, err)
	}

	if s.UserID() != userID {
		return errors.New("invalid session")
	}

	s.Stop(timeNow)
	ld.activeSessions.Delete(sessionID)
	err = ld.db.FinishSession(ctx, userID, sessionID, timeNow)
	if err != nil { // todo: create a pool of unfinished session and finish them asynchronously?
		return fmt.Errorf("finish session %v: %w", sessionID, err)
	}

	return nil
}
