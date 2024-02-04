package token

import "time"

type Maker interface {
	CreateToken(user_id int64, is_admin bool, email string, duration time.Duration) (string, *Payload, error)

	VerifyToken(token string) (*Payload, error)
}
