package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/inmore/gopaste/internal/model"
	"github.com/inmore/gopaste/internal/storage"
	"go.uber.org/zap"
)

type Server struct {
	log *zap.Logger
	st  storage.Storage
}

func New(log *zap.Logger, st storage.Storage) *Server {
	return &Server{log: log, st: st}
}

func (s *Server) Routes() http.Handler {
	r := chi.NewRouter()
	r.Post("/pastes", s.createPaste)
	r.Get("/pastes/{id}", s.getPaste)
	r.Get("/health", s.health)
	return r
}

type createReq struct {
	Content string `json:"content"`
	TTL     int    `json:"ttl_seconds"`
}

// @Summary Create paste
func (s *Server) createPaste(w http.ResponseWriter, r *http.Request) {
	var req createReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	id := uuid.NewString()
	p := &model.Paste{
		ID:        id,
		Content:   req.Content,
		TTL:       req.TTL,
		ExpiresAt: time.Now().Add(time.Duration(req.TTL) * time.Second),
	}
	if err := s.st.Save(p); err != nil {
		http.Error(w, "save error", http.StatusInternalServerError)
		return
	}
	resp := struct {
		ID        string    `json:"id"`
		ExpiresAt time.Time `json:"expires_at"`
	}{id, p.ExpiresAt}
	_ = json.NewEncoder(w).Encode(resp)
}

// @Summary Get paste
func (s *Server) getPaste(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	p, err := s.st.Load(id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	_ = json.NewEncoder(w).Encode(p)
}

// @Summary Health
func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
