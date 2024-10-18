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

// swagger:model signInRequest
// SignInRequest is a struct that defines the request body for the sign-in endpoint.
type SignInRequest struct {
	// Username of the user.
	//
	// required: true
	// example: meliponeech
	Username string `json:"username"`

	// Password of the user.
	//
	// required: true
	// example: whoAmIToJudgey6458yt875yt8754ytierytiueuy
	Password string `json:"password"`
}

// SignInResponse is a struct that defines the response body for the sign-in endpoint.
type SignInResponse struct {
	// Authorization token.
	//
	// example: Bearer rhdfiugreuherugtherutgherutheruthyeruthyerutheuty478t8475t87845y8574y745ty57s
	Authorization string `json:"authorization"`
}

// GetUserInfoRequest is a struct that defines the request body for the getUserInfo endpoint.
type GetUserInfoRequest struct {
}

// swagger:model getUserInfoResponse
// GetUserInfoResponse is a struct that defines the response body for the getUserInfo endpoint.
type GetUserInfoResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// GenericError is a generic error message returned by a server.
type GenericError struct {
	// The error message.
	//
	// example: No no no mister fish you won't go into tazik
	Message string `json:"message"`
}
