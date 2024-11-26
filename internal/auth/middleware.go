package auth

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

const USER_ID_KEY = "user_id"

func (s *Server) authMiddleware(next http.HandlerFunc) http.HandlerFunc	{
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		userID, err := s.validateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), USER_ID_KEY, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}	
}

func userIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(USER_ID_KEY).(uuid.UUID)
	return userID, ok
}
