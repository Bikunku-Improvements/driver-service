package common

import (
	"encoding/json"
	"fmt"
	"github.com/TA-Aplikasi-Pengiriman-Barang/driver-service/internal/domain"
	"github.com/golang-jwt/jwt"
	"os"
	"time"
)

type MyClaims struct {
	jwt.StandardClaims
	Data domain.Bus `json:"data"`
}

func NewJWT(bus domain.Bus) (string, error) {
	claims := MyClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "driver-service",
			ExpiresAt: time.Now().Add(time.Hour * 8).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Data: bus,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func parseJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func ExtractTokenData(tokenString string) (*MyClaims, error) {
	token, err := parseJWT(tokenString)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("unauthorized")
	}

	data := claims["data"]
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var bus domain.Bus
	err = bus.UnmarshalBinary(b)
	if err != nil {
		return nil, err
	}

	return &MyClaims{
		Data: bus,
	}, nil
}
