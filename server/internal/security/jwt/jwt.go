package jwtutil

import (
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type claims struct {
	UserID uint `json:"uid"`
	jwt.RegisteredClaims
}

func Sign(userID uint, secret string, ttl time.Duration) (string, error) {
	now := time.Now()
	c := claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, &c)
	return t.SignedString([]byte(secret))
}

func ParseUserID(tokenString, secret string) (uint, error) {
	var c claims
	token, err := jwt.ParseWithClaims(tokenString, &c, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return 0, err
	}
	if !token.Valid {
		return 0, errors.New("invalid token")
	}
	return c.UserID, nil
}
