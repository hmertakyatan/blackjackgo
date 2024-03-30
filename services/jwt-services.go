package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/hmertakyatan/blackjackgo/utils"
)

type JwtService struct{}

func NewJwtService() *JwtService {
	return &JwtService{}
}

func (j *JwtService) GenerateToken(payload interface{}) (string, error) {
	jwtconfigs, err := utils.JwtConfigLoader()
	if err != nil {
		return "", err
	}
	td := jwtconfigs.AccessTokenExpireDuration
	SecretJWTKey := jwtconfigs.JwtSercretKey
	token := jwt.New(jwt.SigningMethodHS256)
	now := time.Now().UTC()
	claim := token.Claims.(jwt.MapClaims)
	claim["sub"] = payload
	claim["exp"] = now.Add(td).Unix()
	claim["iat"] = now.Unix()
	claim["nbf"] = now.Unix()

	tokenString, err := token.SignedString([]byte(SecretJWTKey))

	if err != nil {
		fmt.Println("Got an error when generating JWToken. ERROR: ", err)
		return "", err
	}

	return tokenString, nil
}

func (j *JwtService) ValidateToken(token string) error {

	jwtconfigs, err := utils.JwtConfigLoader()
	if err != nil {
		return err
	}
	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {

			return nil, fmt.Errorf("Unexpected method %s", jwtToken.Header["alg"])
		}

		return []byte(jwtconfigs.JwtSercretKey), nil
	})

	if err != nil {
		fmt.Println("Invalid token, ERROR: ", err)
		return err
	}

	_, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		fmt.Println("Invalid token claim ERROR: ", err)
		return err
	}

	return nil
}

func (j *JwtService) ExtractAllClaimsFromToken(tokenString string) (map[string]interface{}, error) {
	jwtconfigs, err := utils.JwtConfigLoader()
	signedJWTKey := jwtconfigs.JwtSercretKey
	if err != nil {
		return nil, err
	}
	token, err := jwt.Parse(tokenString, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected method %s", jwtToken.Header["alg"])
		}
		return []byte(signedJWTKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("Failed to parse JWT token: %s", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("Invalid JWT token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("Failed to extract claims from JWT token")
	}

	return map[string]interface{}(claims), nil
}

func (j *JwtService) ExtractPlayerIdClaimFromToken(tokenString string) (string, error) {
	jwtconfigs, err := utils.JwtConfigLoader()
	signedJWTKey := jwtconfigs.JwtSercretKey
	if err != nil {
		return "", err
	}

	token, err := jwt.Parse(tokenString, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected method %s", jwtToken.Header["alg"])
		}
		return []byte(signedJWTKey), nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("Invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("Invalid token claims")
	}

	playerID, ok := claims["player_id"].(string)
	if !ok {
		return "", errors.New("Player ID claim not found")
	}

	return playerID, nil
}
