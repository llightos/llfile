package service

import (
	"github.com/dgrijalva/jwt-go"

	//"net/http"
	"time"
)

const (
	SecretKey  = "en"
	PrivateKey = "lightos"
)

type JWTKey struct {
	Key []byte
}

type TargetUser struct {
	UserID int    `json:"user_id"`
	Mix    string `json:"user_id的AES混淆"`
	//Password string`json:"password"`
	jwt.StandardClaims
}

// 生成一个Token
func (j *JWTKey) CreateToken(claim TargetUser) (string, error) {
	claim.ExpiresAt = time.Now().Add(2 * time.Hour).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString(j.Key)
}

// 解析一个Token
func (j *JWTKey) ParserToken(tokenString string) (*TargetUser, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TargetUser{}, func(token *jwt.Token) (interface{}, error) {
		return j.Key, nil
	})

	if token != nil {
		//Valid : Is the token valid?  Populated when you Parse/Verify a token
		if claims, ok := token.Claims.(*TargetUser); ok && token.Valid {
			return claims, nil
		}
	}
	return nil, err
}
