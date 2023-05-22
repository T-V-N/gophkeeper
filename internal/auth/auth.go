package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go/v4"

	"github.com/T-V-N/gophkeeper/internal/config"
	"github.com/T-V-N/gophkeeper/internal/utils"
)

type UIDKey struct{}

type Claims struct {
	jwt.StandardClaims
	UID string
}

func CreateToken(uid string, cfg *config.Config) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(time.Duration(cfg.JWTExpireTiming * int64(time.Hour)))),
			IssuedAt:  jwt.At(time.Now()),
		},
		UID: uid,
	})
	return token.SignedString([]byte(cfg.SecretKey))
}

func ParseToken(token string, key []byte) (string, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", utils.ErrNotAuthorized
		}
		return key, nil
	})

	if err != nil {
		return "", utils.ErrNotAuthorized
	}

	if claims, ok := parsedToken.Claims.(*Claims); ok && parsedToken.Valid {
		return claims.UID, nil
	}
	return "", utils.ErrNotAuthorized
}
