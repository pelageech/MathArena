package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/go-chi/chi/v5"

	"github.com/pelageech/matharena/internal/data"
	"github.com/pelageech/matharena/internal/models"
	"github.com/pelageech/matharena/internal/pkg/ioutil"
)

type Authorization struct {
	data   Datalayer
	ew     ErrorWriter
	logger Logger
}

func NewAuthorization(data Datalayer, ew ErrorWriter, logger Logger) *Authorization {
	return &Authorization{
		data:   data,
		ew:     ew,
		logger: logger,
	}
}

// swagger:model signUpRequest
// SignUpRequest is a struct that defines the request body for the sign-up endpoint.
type SignUpRequest struct {
	// Username of the user.
	//
	// required: true
	// example: user123
	Username string `json:"username"`

	// Email of the user.
	//
	// required: true
	// example: user@example.com
	Email string `json:"email"`

	// Password of the user.
	//
	// required: true
	// example: myVerySecurePassword123
	Password string `json:"password"`
}

// swagger:route POST /api/signup SignUp
//
// Creates a new user.
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// Schemes: http
//
// Parameters:
// + name: request
//	 in: body
//   description: Sign up request.
//   required: true
//   type: signUpRequest
//
//
// Responses:
//  201: signUpCreatedResponse
//  400: signUpBadRequestError
//  409: signUpConflictError
//  500: signUpInternalServerError

// SignUp is a handler for the sign-up endpoint.
func (a *Authorization) SignUp(w http.ResponseWriter, r *http.Request) {
	// We always return JSON from our API
	w.Header().Set("Content-Type", "application/json")

	var request SignUpRequest

	// We use the FromJSON function to deserialize the request body
	// because it is faster than using the json.Unmarshal function
	err := ioutil.FromJSON(&request, r.Body)
	if err != nil {
		a.ew.Error(w, "Unable to unmarshal JSON", http.StatusBadRequest)
		return
	}

	if err := validate(request); err != nil {
		a.ew.Error(w, fmt.Sprintf("sign up: %v", err), http.StatusBadRequest)
		return
	}

	// We create separate models for API request and datalayer request
	// because we don't want to expose the datalayer models to the API
	// users. This is a good practice to follow.
	// And also in case we want to change the datalayer models or the
	// API request models, we can do it without affecting the other.
	err = a.data.CreateUser(r.Context(), models.User{
		Username: request.Username,
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		if errors.Is(err, data.ErrEmailOrUsernameExists) {
			a.ew.Error(w, data.ErrEmailOrUsernameExists.Error(), http.StatusConflict)
			return
		}

		a.logger.Error("Unable to create user", "error", err)
		a.ew.Error(w, models.ErrInternalServer.Error(), http.StatusInternalServerError)

		return
	}

	// We set the status code to 201 to indicate that the resource is created
	w.WriteHeader(http.StatusCreated)
}

var _usernameValid = regexp.MustCompile(`^[A-Za-z0-9]+$`)
var _emailValid = regexp.MustCompile(`^[\w-.]+@([\w-]+\.)+[\w-]{2,4}$`)

var _filter = []string{"huy", "pizd", "xuy", "xyu", "pidor",
	"ebl", "gavn", "suka", "manda", "mudak",
	"mydak", "sex", "cekc", "hui", "siski", "jopa",
	"boobs",
}

func verifyUsernameWords(s string) bool {
	for _, f := range _filter {
		if strings.Contains(s, f) {
			return false
		}
	}
	return true
}

func verifyPassword(s string) bool {
	var sevenOrMore, number, upper, special bool
	letters := 0
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
			letters++
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		case unicode.IsLetter(c) || c == ' ':
			letters++
		}
	}
	sevenOrMore = len(s) >= 7
	return sevenOrMore && number && upper && special
}

func validate(req SignUpRequest) error {
	if len(req.Username) < 3 || len(req.Username) > 50 {
		return errors.New("username must contain from 3 to 50 symbols")
	}

	if !verifyUsernameWords(req.Username) {
		return errors.New("username contains forbidden symbols")
	}

	if !_usernameValid.MatchString(req.Username) {
		return errors.New("username must contain latin letters or digits")
	}

	if !_emailValid.MatchString(req.Email) {
		return errors.New("email invalid")
	}

	if !verifyPassword(req.Password) {
		return errors.New("a password must be seven or more characters including one uppercase letter," +
			" one special character and alphanumeric characters")
	}

	return nil
}

// swagger:route POST /api/signin SignIn
// Signs in a user.
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// Schemes: http
//
// Parameters:
// + name: request
//   in: body
//   description: Sign in request.
//	 required: true
//	 type: signInRequest
//
//
// Responses:
// 200: signInOkResponse
// 400: signInBadRequestError
// 401: signInUnauthorizedError
// 500: signInInternalServerError

// SignIn is a handler for the sign-in endpoint.
func (a *Authorization) SignIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request models.SignInRequest

	err := ioutil.FromJSON(&request, r.Body)
	if err != nil {
		a.ew.Error(w, "Unable to unmarshal JSON", http.StatusBadRequest)
		return
	}

	authToken, err := a.data.SignInUser(r.Context(), request.Username, request.Password)
	if err != nil {
		if errors.Is(err, models.ErrUnauthorized) {
			a.ew.Error(w, models.ErrUnauthorized.Error(), http.StatusUnauthorized)
			return
		}

		if errors.Is(err, models.ErrUserNotFound) {
			a.ew.Error(w, models.ErrUserNotFound.Error(), http.StatusUnauthorized)
			return
		}

		a.logger.Error("Unable to get Bearer token", "error", err)
		a.ew.Error(w, models.ErrInternalServer.Error(), http.StatusInternalServerError)

		return
	}

	err = ioutil.ToJSON(models.SignInResponse{
		Authorization: authToken,
	}, w)
	if err != nil {
		// log the error to debug it
		a.logger.Error("Unable to write JSON response", "error", err)
		// write a generic error to the response writer, so we don't expose the actual error
		a.ew.Error(w, models.ErrInternalServer.Error(), http.StatusInternalServerError)

		return
	}
}

// swagger:route GET /api/user/{id} GetUserInfo
// Get user info.
//
// Produces:
// - application/json
//
// Schemes: http
//
// Parameters:
// + name: id
//   in: query
//   description: UserId.
//   required: true
//   type: integer
//
// Responses:
// 200: getUserInfoResponse
// 400: getUserInfoBadRequestError
// 404: getUserInfoNotFoundError
// 500: getUserInfoInternalServerError

// GetUserInfo is a handler for the get-user-info endpoint.
func (a *Authorization) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		a.ew.Error(w, models.ErrInternalServer.Error(), http.StatusInternalServerError)
		a.logger.Error("Unable to convert user id to integer", "error", err)
		return
	}

	if userID < 1 {
		a.ew.Error(w, "User id cannot be less than 1", http.StatusBadRequest)
		return
	}

	userInfo, err := a.data.GetUserById(r.Context(), userID)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			a.ew.Error(w, models.ErrUserNotFound.Error(), http.StatusNotFound)
			return
		}
		a.ew.Error(w, models.ErrInternalServer.Error(), http.StatusInternalServerError)
		a.logger.Error("unable to get user info in GetUserInfo", "error", err, "userId", userID)
		return
	}

	err = ioutil.ToJSON(models.GetUserInfoResponse(userInfo), w)
	if err != nil {
		a.ew.Error(w, models.ErrInternalServer.Error(), http.StatusInternalServerError)
		a.logger.Error("unable to marshal userInfo in GetUserInfo", "error", err, "userInfo", userInfo)
		return
	}
}
