package services

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateToken(td time.Duration, payload interface{}, SecretJWTKey string) (string, error) {
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

func ValidateToken(token string, signedJWTKey string) (interface{}, error) {
	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {

			return nil, fmt.Errorf("Unexpected method %s", jwtToken.Header["alg"])
		}

		return []byte(signedJWTKey), nil
	})

	if err != nil {
		fmt.Println("Invalid token, ERROR: ", err)
		return nil, err
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		fmt.Println("Invalid token claim ERROR: ", err)
		return nil, err
	}

	return claims["sub"], nil
}
