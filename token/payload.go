package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrorInvalidToken = errors.New("token is invalid")
	ErrorExpiredToken = errors.New("token has expired")
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	UserId    int64     `json:"user_id"`
	Email     string    `json:"email"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(user_id int64, email string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenID,
		UserId:    user_id,
		Email:     email,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrorExpiredToken
	}
	return nil
}
