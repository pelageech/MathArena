package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrUserNotFound = fmt.Errorf("user not found")

// GetUserInfo returns username, email and error by given userId.
func (d *PSQLDatabase) GetUserInfo(ctx context.Context, userId int) (username string, email string, err error) {
	row := d.QueryRow(ctx, `
SELECT username, email FROM players WHERE id = $1
`,
		userId)

	if err := row.Scan(&username, &email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", ErrUserNotFound
		}

		return "", "", fmt.Errorf("unable to get username and email from db in GetUserInfo: %w", err)
	}

	return username, email, nil
}

// HasEmailOrUsername checks whether a user with the given username or email exists.
func (d *PSQLDatabase) HasEmailOrUsername(ctx context.Context, username, email string) (bool, error) {
	var count int

	row := d.QueryRow(ctx, `
SELECT COUNT(*) FROM players WHERE username = $1 OR email = $2
`,
		username, email)

	if err := row.Scan(&count); err != nil {
		return false, fmt.Errorf("unable to count rows in HasEmailOrUsername: %w", err)
	}

	if count == 0 {
		return false, nil
	}
	return true, nil
}

// InsertUser inserts a new user into the database.
func (d *PSQLDatabase) InsertUser(ctx context.Context, username, hashedPassword, email string) (id int, err error) {
	row := d.QueryRow(ctx, `
INSERT INTO players (email, username, hashed_password)
VALUES ($1, $2, $3) RETURNING id
`,
		email, username, hashedPassword)

	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("unable to insert user in InsertUser: %w", err)
	}

	return id, nil
}

// GetHashedPassword returns salt and hash for the given username.
func (d *PSQLDatabase) GetHashedPassword(ctx context.Context, username string) (hashedPassword string, err error) {
	row := d.QueryRow(ctx, `
SELECT hashed_password FROM players WHERE username = $1
`,
		username)

	if err := row.Scan(&hashedPassword); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrUserNotFound
		}

		return "", fmt.Errorf("unable to get salt and hash in GetHashedPassword: %w", err)
	}

	return hashedPassword, nil
}

// GetUserID returns user id for the given username.
func (d *PSQLDatabase) GetUserID(ctx context.Context, username string) (userId int64, err error) {
	row := d.QueryRow(ctx, `
SELECT id FROM players WHERE username = $1
`,
		username)

	if err := row.Scan(&userId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrUserNotFound
		}

		return 0, fmt.Errorf("unable to get user id in GetUserID: %w", err)
	}

	return userId, nil
}

func (d *PSQLDatabase) GetAllIds(ctx context.Context) (ids []int, err error) {
	rows, err := d.Query(ctx, `
SELECT id FROM players
`)
	defer rows.Close()

	if err != nil {
		d.Logger().Errorf("unable to get all user credentials: %v", err)
		return nil, fmt.Errorf("unable to get all user credentials: %w", err)
	}

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			d.Logger().Errorf("unable to get all user credentials: %v", err)
			return nil, fmt.Errorf("unable to get all user credentials: %w", err)
		}
		ids = append(ids, id)
	}

	return ids, nil
}
