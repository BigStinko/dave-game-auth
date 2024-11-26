package auth

import (
	"encoding/json"
	"net/http"

	"aidanwoods.dev/go-paseto"
	"github.com/BigStinko/dave-game-auth/internal/db"
	"golang.org/x/crypto/bcrypt"
)

type Server struct {
	db *db.Queries
	symmetricKey paseto.V4SymmetricKey
}

type RegisterRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type AuthResponse struct {
    Token string `json:"token"`
    Error string `json:"error,omitempty"`
}

func NewServer(db *db.Queries, key paseto.V4SymmetricKey) *Server {
	return &Server{
		db: db,
		symmetricKey: key,
	}
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	user, err := s.db.CreateUser(r.Context(), db.CreateUserParams{
		Username: req.Username,
		PasswordHash: string(hashedPassword),
	})
	if err != nil {
		http.Error(w, "Username already taken", http.StatusConflict)
		return
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(AuthResponse{Token: token})
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := s.db.GetUserByUsername(r.Context(), req.Username)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(AuthResponse{Token: token})
}

func (s *Server) handleGetMatchHistory(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	matches, err := s.db.GetUserMatches(r.Context(), userID)
	if err != nil {
		http.Error(w, "Error fetching matches", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(matches)
}

func (s *Server) handleGetStats(w http.ResponseWriter, r *http.Request) {
	userId, ok := userIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	stats, err := s.db.GetUserStats(r.Context(), userId)
	if err != nil {
		http.Error(w, "Error fetching stats", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(stats)
}
