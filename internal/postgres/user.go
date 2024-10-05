package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrUserNotFound = fmt.Errorf("user not found")

// GetUserInfo returns username, email and error by given userId.
func (d *PSQLDatabase) GetUserInfo(ctx context.Context, userId int) (string, string, error) {
	var username string
	var email string

	err := d.QueryRow(ctx, "SELECT username, email FROM user_credentials WHERE user_id = $1", userId).Scan(&username, &email)
	if err != nil {
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

	err := d.QueryRow(ctx, "SELECT COUNT(*) FROM user_credentials WHERE username = $1 OR email = $2", username, email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("unable to count rows in HasEmailOrUsername: %w", err)
	}

	if count == 0 {
		return false, nil
	}
	return true, nil
}

// InsertUser inserts a new user into the database.
func (d *PSQLDatabase) InsertUser(ctx context.Context, username, salt, hash, email string) (int, error) {
	row := d.QueryRow(ctx, "INSERT INTO user_credentials (username, salt, hash, email) VALUES ($1, $2, $3, $4) RETURNING user_id", username, salt, hash, email)
	var id int

	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("unable to insert user in InsertUser: %w", err)
	}

	return id, nil
}

// GetSaltAndHash returns salt and hash for the given username.
func (d *PSQLDatabase) GetSaltAndHash(ctx context.Context, username string) (salt, hash string, err error) {
	err = d.QueryRow(ctx, "SELECT salt, hash FROM user_credentials WHERE username = $1", username).Scan(&salt, &hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", ErrUserNotFound
		}

		return "", "", fmt.Errorf("unable to get salt and hash in GetSaltAndHash: %w", err)
	}

	return salt, hash, nil
}

// GetUserID returns user id for the given username.
func (d *PSQLDatabase) GetUserID(ctx context.Context, username string) (int64, error) {
	var userID int64

	err := d.QueryRow(ctx, "SELECT user_id FROM user_credentials WHERE username = $1", username).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrUserNotFound
		}

		return 0, fmt.Errorf("unable to get user id in GetUserID: %w", err)
	}

	return userID, nil
}

func (d *PSQLDatabase) GetAllIds(ctx context.Context) ([]int, error) {
	rows, err := d.Query(ctx, "SELECT user_id FROM user_credentials")
	if err != nil {
		d.Logger().Errorf("unable to get all user credentials: %v", err)
		return nil, fmt.Errorf("unable to get all user credentials: %w", err)
	}
	defer rows.Close()
	var ids []int
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
