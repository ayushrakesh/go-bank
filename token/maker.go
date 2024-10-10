package token

import "time"

type Maker interface {
	CreateToken(string, time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}
