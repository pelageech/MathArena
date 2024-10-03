package models

// User is a struct that defines the user model.
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Users is a struct that defines a slice of User.
type Users []User

// AuthorizationToken is a type that defines the authorization token.
type AuthorizationToken string

// UserInfo is a struct that defines the user info model.
type UserInfo struct {
	ID       int
	Username string
	Email    string
}

// SignInRequest is a struct that defines the request body for the sign-in endpoint.
type SignInRequest struct {
	// Username of the user.
	Username string `json:"username"`

	// Password of the user.
	Password string `json:"password"`
}

// SignInResponse is a struct that defines the response body for the sign-in endpoint.
type SignInResponse struct {
	// Authorization token.
	Authorization string `json:"authorization"`
}

// GetUserInfoRequest is a struct that defines the request body for the getUserInfo endpoint.
type GetUserInfoRequest struct {
}

// GetUserInfoResponse is a struct that defines the response body for the getUserInfo endpoint.
type GetUserInfoResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// GenericError is a generic error message returned by a server.
type GenericError struct {
	// The error message.
	Message string `json:"message"`
}
