package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type JwtConfig struct {
	AccessTokenExpireDuration time.Duration
	JwtSercretKey             string
}

func JwtConfigLoader() (*JwtConfig, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return nil, err
	}
	jwtSercet := os.Getenv("JWTSecretKey")
	accessTokenExpireStr := os.Getenv("Access_token_expire")
	accessTokenTd, err := time.ParseDuration(accessTokenExpireStr)
	if err != nil {

		return nil, fmt.Errorf("acces token cannot converting")
	}
	jwtconfigs := &JwtConfig{
		JwtSercretKey:             jwtSercet,
		AccessTokenExpireDuration: accessTokenTd,
	}
	return jwtconfigs, nil
}
