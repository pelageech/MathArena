package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/pelageech/matharena/internal/game"
	"time"
)

func (p *PSQLDatabase) CreateSession(ctx context.Context, userId int, startTime time.Time) (game.SessionID, error) {
	_, _, err := p.GetUserInfo(ctx, userId)
	if err != nil {
		return -1, fmt.Errorf("%v: %w", userId, ErrUserNotFound)
	}

	row := p.QueryRow(ctx, `INSERT INTO game_sessions(player_id, start_time, end_time, points, is_finished) VALUES 
                              ($1, $2, to_timestamp(0), 0, false) RETURNING id`,
		userId,
		startTime,
	)
	var id int
	err = row.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error creating session: %w", err)
	}

	return game.SessionID(id), nil
}
func (p *PSQLDatabase) FinishSession(ctx context.Context, userID int, id game.SessionID, finishTime time.Time) error {
	actualUserID, err := p.GetUserIDBySession(ctx, id)
	if err != nil {
		return err
	}

	if actualUserID != userID {
		return errors.New("forbidden operation")
	}

	_, err = p.Exec(ctx, `UPDATE game_sessions SET is_finished = true, end_time = $1 WHERE id = $2`,
		finishTime,
		id,
	)
	if err != nil {
		return fmt.Errorf("error updating session: %w", err)
	}

	return nil
}

func (p *PSQLDatabase) GetUserIDBySession(ctx context.Context, id game.SessionID) (int, error) {
	row := p.QueryRow(ctx, `SELECT player_id FROM game_sessions WHERE id = $1`,
		id,
	)

	var gotID int
	err := row.Scan(&gotID)
	if err != nil {
		return 0, fmt.Errorf("error getting user id: %w", err)
	}

	return gotID, nil
}
