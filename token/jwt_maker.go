package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

type JwtMaker struct {
	secretKey string
}

func NewJwtMaker(secretKey string) (Maker, error) {
	return &JwtMaker{secretKey}, nil
}

func (maker *JwtMaker) CreateToken(user_id int64, email string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(user_id, email, duration)
	if err != nil {
		return "", payload, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(maker.secretKey))
	return token, payload, err
}

func (maker *JwtMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrorInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrorExpiredToken) {
			return nil, ErrorExpiredToken
		}
		return nil, ErrorInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrorInvalidToken
	}

	return payload, nil
}