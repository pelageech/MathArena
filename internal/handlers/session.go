package handlers

import (
	"context"
	"github.com/charmbracelet/log"
	"github.com/pelageech/matharena/internal/game"
	"github.com/pelageech/matharena/internal/game/generator"
	"github.com/pelageech/matharena/internal/pkg/ioutil"
	"net/http"
	"time"
)

type GameSessionsDatalayer interface {
	CreateSession(context.Context, time.Duration, int, generator.Difficulty, time.Time) (*game.Session, error)
	Answer(context.Context, game.SessionID, int, int, time.Time) (*game.Session, error)
	Stop(context.Context, game.SessionID, int, time.Time) error
}

type GameSessionsHandler struct {
	data   GameSessionsDatalayer
	ew     ErrorWriter
	logger *log.Logger
}

func NewGameSessionsHandler(data GameSessionsDatalayer, ew ErrorWriter, logger *log.Logger) *GameSessionsHandler {
	return &GameSessionsHandler{data: data, ew: ew, logger: logger}
}

type CreateSessionRequest struct {
	UserID int `json:"user_id"`
}

type CreateSessionResponse struct {
	SessionID  string        `json:"session_id"`
	TimeLeft   time.Duration `json:"time_left"`
	Expression string        `json:"expression"`
	Score      int           `json:"score"`
}

func (h *GameSessionsHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")

	reqBody := CreateSessionRequest{}
	if err := ioutil.FromJSON(&reqBody, r.Body); err != nil {
		h.ew.Error(w, "Unable to unmarshal JSON", http.StatusBadRequest)
		return
	}

	s, err := h.data.CreateSession(r.Context(), time.Minute, reqBody.UserID, generator.Easy, time.Now())
	if err != nil {
		h.logger.Errorf("unable to create session: %v", err)
		h.ew.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respBody := CreateSessionResponse{
		SessionID:  s.ID().String(),
		TimeLeft:   s.TimeLeft(),
		Expression: string(s.CurrentExpression().Marshal()),
		Score:      s.Score(),
	}
	if err := ioutil.ToJSON(respBody, w); err != nil {
		h.logger.Errorf("unable to marshal response: %v", err)
		h.ew.Error(w, "unable to marshal JSON", http.StatusInternalServerError)
		return
	}
}

type AnswerRequest struct {
	SessionID string `json:"session_id"`
	UserID    int    `json:"user_id"`
	Answer    int    `json:"answer"`
}

type AnswerResponse struct {
	SessionID  string        `json:"session_id"`
	TimeLeft   time.Duration `json:"time_left"`
	Expression string        `json:"expression"`
	Score      int           `json:"score"`
}

func (h *GameSessionsHandler) Answer(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")

	reqBody := AnswerRequest{}
	if err := ioutil.FromJSON(&reqBody, r.Body); err != nil {
		h.ew.Error(w, "Unable to unmarshal JSON", http.StatusBadRequest)
		return
	}

	id, err := game.ParseSessionID(reqBody.SessionID)
	if err != nil {
		h.ew.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s, err := h.data.Answer(r.Context(), id, reqBody.Answer, reqBody.UserID, time.Now())
	if err != nil {
		h.logger.Errorf("unable to answer: %v", err)
		h.ew.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respBody := AnswerResponse{
		SessionID:  s.ID().String(),
		TimeLeft:   s.TimeLeft(),
		Expression: string(s.CurrentExpression().Marshal()),
		Score:      s.Score(),
	}
	if err := ioutil.ToJSON(respBody, w); err != nil {
		h.logger.Errorf("unable to marshal response: %v", err)
		h.ew.Error(w, "unable to marshal JSON", http.StatusInternalServerError)
		return
	}
}

type StopRequest struct {
	SessionID string `json:"session_id"`
	UserID    int    `json:"user_id"`
}

func (h *GameSessionsHandler) Stop(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")

	reqBody := StopRequest{}
	if err := ioutil.FromJSON(&reqBody, r.Body); err != nil {
		h.ew.Error(w, "Unable to unmarshal JSON", http.StatusBadRequest)
		return
	}

	id, err := game.ParseSessionID(reqBody.SessionID)
	if err != nil {
		h.ew.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.data.Stop(r.Context(), id, reqBody.UserID, time.Now())
	if err != nil {
		h.logger.Errorf("unable to stop: %v", err)
		h.ew.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
