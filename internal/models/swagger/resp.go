// Package swagger describes models for application.
//
// Documentation for MathArena API.
//
// Schemes: http
// BasePath: /
// Version: 0.0.1
// Contact: Artyom Blaginin<pelageech@mail.ru>
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// swagger:meta
package swagger

import "github.com/pelageech/matharena/internal/models"

// swagger:response signUpCreatedResponse
type SignUpCreatedResponse struct{}

// swagger:response signUpBadRequestError
type SignUpBadRequestError struct {
	// in: body
	Body models.GenericError
}

// swagger:response signUpConflictError
type SignUpConflictError struct {
	// in: body
	Body models.GenericError
}

// swagger:response signUpInternalServerError
type InternalServerError struct {
	// in: body
	Body models.GenericError
}

// swagger:response signInOkResponse
type SignInOkResponse struct {
	// Response with the authorization token.
	// in: body
	Body models.SignInResponse
}

// swagger:response signInBadRequestError
type SignInBadRequestError struct {
	// in: body
	Body models.GenericError
}

// swagger:response signInUnauthorizedError
type SignInUnauthorizedError struct {
	// in: body
	Body models.GenericError
}

// swagger:response signInInternalServerError
type SignInInternalServerError struct {
	// in: body
	Body models.GenericError
}

// swagger:response getUserInfoBadRequestError
type GetUserInfoBadRequestError struct {
	// in: body
	Body models.GenericError
}

// swagger:response getUserInfoNotFoundError
type GetUserInfoNotFoundError struct {
	// in: body
	Body models.GenericError
}

// swagger:response getUserInfoInternalServerError
type GetUserInfoInternalServerError struct {
	// in: body
	Body models.GenericError
}
