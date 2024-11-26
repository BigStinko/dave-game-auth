package auth

import (
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/google/uuid"
)

func (s *Server) generateToken(userID uuid.UUID) (string, error) {
	token := paseto.NewToken()
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(24 * time.Hour))
	token.Set(USER_ID_KEY, userID)
	return token.V4Encrypt(s.symmetricKey, nil), nil
}

func (s *Server) validateToken(tokenString string) (uuid.UUID, error) {
	parser := paseto.NewParser()
	parser.AddRule(paseto.NotExpired())
	parser.AddRule(paseto.ValidAt(time.Now()))

	token, err := parser.ParseV4Local(s.symmetricKey, tokenString, nil)
	if err != nil {
		return uuid.Nil, err
	}

	userIDStr, err := token.GetString(USER_ID_KEY)
	if err != nil {
		return uuid.Nil, err
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}
