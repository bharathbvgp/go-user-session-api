package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("jwt_secret")

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

func GenerateToken(userID uint) (string , error) {
	expirationTime := time.Now().Add(12 * time.Hour)
	claims := Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256 , claims)
	return token.SignedString(jwtKey)
}

func ValidateToken(tokenString string) (*Claims , error) {
	claims := &Claims{}
	token , err := jwt.ParseWithClaims(tokenString , claims , func (token *jwt.Token) (interface{} , error){
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil , err
		}
		return nil , err
	}
	if !token.Valid {
		return nil,err
	}
	return claims , nil
}
