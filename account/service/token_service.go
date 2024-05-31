package service

import (
	"context"
	"crypto/rsa"
	"log"

	"github.com/nobitayon/memorization-app/account/handler/model"
	"github.com/nobitayon/memorization-app/account/handler/model/apperrors"
)

type tokenService struct {
	TokenRepository       model.TokenRepository
	PrivKey               *rsa.PrivateKey
	PubKey                *rsa.PublicKey
	RefreshSecret         string
	IDExpirationSecs      int64
	RefreshExpirationSecs int64
}

type TSConfig struct {
	TokenRepository       model.TokenRepository
	PrivKey               *rsa.PrivateKey
	PubKey                *rsa.PublicKey
	RefreshSecret         string
	IDExpirationSecs      int64
	RefreshExpirationSecs int64
}

func NewTokenService(c *TSConfig) model.TokenService {
	return &tokenService{
		TokenRepository:       c.TokenRepository,
		PrivKey:               c.PrivKey,
		PubKey:                c.PubKey,
		RefreshSecret:         c.RefreshSecret,
		IDExpirationSecs:      c.IDExpirationSecs,
		RefreshExpirationSecs: c.RefreshExpirationSecs,
	}
}

func (s *tokenService) NewPairFromUser(ctx context.Context, u *model.User, prevTokenID string) (*model.TokenPair, error) {
	idToken, err := generateIDToken(u, s.PrivKey, s.IDExpirationSecs)
	if err != nil {
		log.Printf("error generating idToken for uid: %v. Error: %v", u.UID, err.Error())
		return nil, apperrors.NewInternal()
	}

	refreshToken, err := GenerateRefreshToken(u.UID, s.RefreshSecret, s.RefreshExpirationSecs)
	if err != nil {
		log.Printf("error generating refresh token for uid: %v. Error: %v\n", u.UID, err.Error())
		return nil, apperrors.NewInternal()
	}

	// todo: store refresh token by calling TokenRepository methods
	if err := s.TokenRepository.SetRefreshToken(ctx, u.UID.String(), refreshToken.ID,
		refreshToken.ExpiresIn); err != nil {
		log.Printf("error storing tokenID for uid: %v. Error: %v\n", u.UID, err.Error())
		return nil, apperrors.NewInternal()
	}

	// delete user's current refresh token (used when refreshing idToken)
	if prevTokenID != "" {
		if err := s.TokenRepository.DeleteRefreshToken(ctx, u.UID.String(), prevTokenID); err != nil {
			log.Printf("could not delete previous refreshToken for uid: %v, tokenID: %v", u.UID.String(), prevTokenID)
		}
	}

	return &model.TokenPair{
		IDToken:      idToken,
		RefreshToken: refreshToken.SS,
	}, nil
}
