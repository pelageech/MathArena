package data

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/pelageech/matharena/internal/models"
)

// ErrEmailOrUsernameExists is an error returned when a user with the given email already exists.
var ErrEmailOrUsernameExists = fmt.Errorf("user with this email or username already exists")

// UserCredentials is an interface that represents database.
//
//go:generate mockery --name UserCredentials --output=./ --filename=mocks/userCredentials.go --with-expecter
type UserCredentials interface {
	HasEmailOrUsername(ctx context.Context, username, email string) (bool, error)
	GetSaltAndHash(ctx context.Context, username string) (salt, hash string, err error)
	InsertUser(ctx context.Context, username, salt, hash, email string) (int, error)
	GetUserInfo(ctx context.Context, userId int) (username, email string, err error)
	GetUserID(ctx context.Context, username string) (int64, error)
}

// Datalayer is a struct that helps us to interact with the data.
type Datalayer struct {
	// db is a database.
	db UserCredentials

	// saltLength is the length of the salt.
	// We will use this to generate a salt for the password.
	saltLength int

	// tokenExpirationTime is the time after which the token will expire.
	tokenExpirationTime time.Duration

	// signKey is a key to sign the token.
	signKey []byte
}

func (d *Datalayer) GetUserById(ctx context.Context, id int) (models.UserInfo, error) {
	username, email, err := d.db.GetUserInfo(ctx, id)
	if err != nil {
		return models.UserInfo{}, fmt.Errorf("unable to get user info in GetUserById: %w", err)
	}

	return models.UserInfo{
		ID:       id,
		Username: username,
		Email:    email,
	}, nil
}

// New returns a new Datalayer struct.
func New(db UserCredentials, saltLength int, tokenExpirationTime time.Duration, signKey []byte) *Datalayer {
	return &Datalayer{
		db:                  db,
		saltLength:          saltLength,
		tokenExpirationTime: tokenExpirationTime,
		signKey:             signKey,
	}
}

// isEmailOrUsernameExists is a helper function to check if a user with the given email already exists.
// Returns true if the user exists, false otherwise.
func (d *Datalayer) isEmailOrUsernameExists(ctx context.Context, username, email string) (bool, error) {
	has, err := d.db.HasEmailOrUsername(ctx, username, email)
	if err != nil {
		return false, fmt.Errorf("unable to check if email exists in isEmailOrUsernameExists: %w", err)
	}

	return has, nil
}

// isPasswordCorrect is a helper function to check if the given password is correct for the given email.
func (d *Datalayer) isPasswordCorrect(ctx context.Context, username, password string) (bool, error) {
	salt, hash, err := d.db.GetSaltAndHash(ctx, username)
	if err != nil {
		return false, fmt.Errorf("unable to get salt and hash in isPasswordCorrect: %w", err)
	}

	// Compare the password with the hash
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password+salt))
	if err != nil {
		return false, nil
	}

	return true, nil
}

func (d *Datalayer) generateSalt() (string, error) {
	saltBytes := make([]byte, d.saltLength)

	_, err := rand.Read(saltBytes)
	if err != nil {
		return "", err
	}

	salt := base64.URLEncoding.EncodeToString(saltBytes)
	return salt, nil
}

// generateToken is a helper function to generate a JWT token.
func (d *Datalayer) generateToken(username string, userID int64) (string, error) {
	// Set the expiration time of the token
	expirationTime := time.Now().Add(d.tokenExpirationTime * time.Minute)

	// Create the JWT claims, which includes the username and expiry time
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"username": username,
		"user_id":  userID,
		"exp":      expirationTime.Unix(),
	})

	// Sign the token with our secret key
	signedToken, err := token.SignedString(d.signKey)
	if err != nil {
		return "", fmt.Errorf("unable to sign token in generateToken: %w", err)
	}

	return signedToken, nil
}

// CreateUser is a function to create a new user.
func (d *Datalayer) CreateUser(ctx context.Context, user models.User) error {
	// Check if user with this email already exists
	has, err := d.isEmailOrUsernameExists(ctx, user.Username, user.Email)
	if err != nil {
		return fmt.Errorf("unable to check if email or username exists in CreateUser: %w", err)
	}

	// If user with this email or username already exists, return an error
	if has {
		return ErrEmailOrUsernameExists
	}

	// Generate salt
	salt, err := d.generateSalt()
	if err != nil {
		return fmt.Errorf("unable to generate salt in CreateUser: %w", err)
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password+salt), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("unable to hash password in CreateUser: %w", err)
	}

	_, err = d.db.InsertUser(ctx, user.Username, salt, string(hash), user.Email)
	if err != nil {
		return fmt.Errorf("unable to insert user in CreateUser: %w", err)
	}

	return nil
}

// SignInUser is a function to sign in a user.
// It returns an authorization token if the user is signed in successfully.
func (d *Datalayer) SignInUser(ctx context.Context, username, password string) (string, error) {
	correct, err := d.isPasswordCorrect(ctx, username, password)
	if err != nil {
		return "", fmt.Errorf("unable to check if password is correct in SignInUser: %w", err)
	}

	if !correct {
		return "", models.ErrUnauthorized
	}

	userID, err := d.db.GetUserID(ctx, username)
	if err != nil {
		return "", fmt.Errorf("unable to get user id in SignInUser: %w", err)
	}

	token, err := d.generateToken(username, userID)
	if err != nil {
		return "", fmt.Errorf("unable to generate token in SignInUser: %w", err)
	}

	return "Bearer " + token, nil
}
