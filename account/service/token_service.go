package service

import (
	"context"
	"crypto/rsa"
	"log"

	"github.com/nobitayon/memorization-app/account/handler/model"
	"github.com/nobitayon/memorization-app/account/handler/model/apperrors"
)

type tokenService struct {
	PrivKey       *rsa.PrivateKey
	PubKey        *rsa.PublicKey
	RefreshSecret string
}

type TSConfig struct {
	PrivKey       *rsa.PrivateKey
	PubKey        *rsa.PublicKey
	RefreshSecret string
}

func NewTokenService(c *TSConfig) model.TokenService {
	return &tokenService{
		PrivKey:       c.PrivKey,
		PubKey:        c.PubKey,
		RefreshSecret: c.RefreshSecret,
	}
}

func (s *tokenService) NewPairFromUser(ctx context.Context, u *model.User, prevTokenID string) (*model.TokenPair, error) {
	idToken, err := generateIDToken(u, s.PrivKey)
	if err != nil {
		log.Printf("error generating idToken for uid: %v. Error: %v", u.UID, err.Error())
		return nil, apperrors.NewInternal()
	}

	refreshToken, err := GenerateRefreshToken(u.UID, s.RefreshSecret)
	if err != nil {
		log.Printf("error generating refresh token for uid: %v. Error: %v\n", u.UID, err.Error())
		return nil, apperrors.NewInternal()
	}

	// todo: store refresh token by calling TokenRepository methods

	return &model.TokenPair{
		IDToken:      idToken,
		RefreshToken: refreshToken.SS,
	}, nil
}
