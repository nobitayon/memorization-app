package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/nobitayon/memorization-app/account/handler"
	"github.com/nobitayon/memorization-app/account/repository"
	"github.com/nobitayon/memorization-app/account/service"
)

func inject(d *dataSource) (*gin.Engine, error) {
	log.Println("injecting data sources")

	// repository layer
	userRepository := repository.NewUserRepository(d.DB)

	// service layer
	userService := service.NewUserService(&service.USConfig{
		UserRepository: userRepository,
	})

	privKeyFile := os.Getenv("PRIV_KEY_FILE")
	priv, err := ioutil.ReadFile(privKeyFile)
	if err != nil {
		return nil, fmt.Errorf("could not read private key pem file: %w", err)
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(priv)
	if err != nil {
		return nil, fmt.Errorf("could not parse private key: %n", err)
	}

	pubKeyFile := os.Getenv("PUB_KEY_FILE")
	pub, err := ioutil.ReadFile(pubKeyFile)
	if err != nil {
		return nil, fmt.Errorf("could not read public key pem file: %w", err)
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pub)
	if err != nil {
		return nil, fmt.Errorf("could not parse public key: %n", err)
	}

	refreshSecret := os.Getenv("REFRESH_SECRET")

	tokenService := service.NewTokenService(&service.TSConfig{
		PrivKey:       privKey,
		PubKey:        pubKey,
		RefreshSecret: refreshSecret,
	})

	router := gin.Default()
	handler.NewHandler(&handler.Config{
		R:            router,
		UserService:  userService,
		TokenService: tokenService,
	})

	return router, nil

}
