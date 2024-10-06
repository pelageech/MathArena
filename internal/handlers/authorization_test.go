package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/charmbracelet/log"
	"github.com/stretchr/testify/mock"

	"github.com/pelageech/matharena/internal/data"
	"github.com/pelageech/matharena/internal/data/mocks"
	"github.com/pelageech/matharena/internal/pkg/ioutil"
)

func TestAuthorization(t *testing.T) {
	cred := mocks.NewUserCredentials(t)

	dl := data.New(cred, 24, time.Hour, []byte{0})

	l := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
	})

	ew := ioutil.JSONErrorWriter{Logger: l}

	authHandlers := NewAuthorization(dl, ew, l)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /signup", authHandlers.SignUp)
	mux.HandleFunc("GET /user/{id}", authHandlers.GetUserInfo)
	mux.HandleFunc("POST /signin", authHandlers.SignIn)

	t.Run("signup", func(t *testing.T) {
		b, _ := json.Marshal(SignUpRequest{
			Username: "aboba",
			Email:    "aboba@g.nsu.ru",
			Password: "Aboba20!8",
		})

		cred.EXPECT().
			HasEmailOrUsername(mock.Anything, "aboba", "aboba@g.nsu.ru").
			Return(false, nil)

		cred.EXPECT().
			InsertUser(mock.Anything, "aboba", mock.Anything, "aboba@g.nsu.ru").
			Return(0, nil)

		req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(b))
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusCreated {
			t.Fatalf("got %d, want %d", res.StatusCode, http.StatusCreated)
		}
	})
}
