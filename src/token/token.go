package token

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte("efcbdb586880e9f5006bc7261b2e1a6e")

type Claims struct {
	UserId int `json:"user_id"`
	jwt.StandardClaims
}

// https://medium.com/企鵝也懂程式設計/225b377e0f79
func NewToken(userId int, username string, t time.Duration) string {
	now := time.Now()

	claims := Claims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			Audience:  username,
			ExpiresAt: now.Add(t).Unix(),
			IssuedAt:  now.Unix(),
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	if err != nil {
		panic(err)
	}

	return token
}

func PareToken(token string) (*Claims, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (i interface{}, e error) {
		return []byte(jwtSecret), nil
	})

	if err == nil && jwtToken != nil {
		if claim, ok := jwtToken.Claims.(*Claims); ok && jwtToken.Valid {
			return claim, nil
		}
	}
	return nil, err
}
