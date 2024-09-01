package jwthandler

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hilmiikhsan/shopeefun-user-service/internal/infrastructure/config"
	"github.com/rs/zerolog/log"
)

func GenerateTokenString(payload CostumClaimsPayload) (string, error) {
	claims := CustomClaims{
		UserId: payload.UserId,
		Role:   payload.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "user",
			Issuer:    "shopeefun-app",
			ExpiresAt: jwt.NewNumericDate(payload.TokenExpiration),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			NotBefore: jwt.NewNumericDate(time.Now().UTC()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	tokenString, err := token.SignedString([]byte(config.Envs.Guard.JwtPrivateKey))
	if err != nil {
		log.Error().Err(err).Msg("jwthandler::GenerateTokenString - Error while signing token")
		return "", err
	}

	return tokenString, nil
}
