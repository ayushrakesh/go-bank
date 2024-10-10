package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PastoMaker struct {
	symmetricKey []byte
	pasto        *paseto.V2
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size, must be exactly of %d size", chacha20poly1305.KeySize)
	}

	maker := &PastoMaker{
		symmetricKey: []byte(symmetricKey),
		pasto:        paseto.NewV2(),
	}
	return maker, nil
}

func (maker *PastoMaker) CreateToken(username string, duration time.Duration) (string, error) {

	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	return maker.pasto.Encrypt(maker.symmetricKey, payload, nil)
}
func (maker *PastoMaker) VerifyToken(token string) (*Payload, error) {

	payload := &Payload{}

	err := maker.pasto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
