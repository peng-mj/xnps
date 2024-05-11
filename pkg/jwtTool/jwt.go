package jwtTool

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"math/rand"
	"time"
)

const (
	MaxRandKeys  = 100
	MaxKeyLength = 64
)

func randByte(l int) []byte {
	base := []byte("0123456789abcdefghijklmnopqrstuvwxyzQWERTYUIOPASDFGHJKLZXCVBNM")
	var res []byte
	for i := 0; i < l; i++ {
		res = append(res, base[rand.Intn(len(base))])
	}
	return res
}

type Token struct {
	keys [][]byte
}

func NewToken() *Token {
	token := Token{keys: make([][]byte, MaxRandKeys)}
	for i := 0; i < MaxRandKeys; i++ {
		token.keys[i] = randByte(MaxKeyLength)
	}
	return &token
}

func (t *Token) Generate(uid int64, expireTime time.Duration) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	sid := rand.Intn(MaxRandKeys)
	claims["uid"] = uid
	claims["sid"] = sid
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(expireTime).Unix()
	tokenString, _ := token.SignedString(t.keys[sid])
	return tokenString
}

func (t *Token) Parse(str string) (uid string, err error) {
	token, err := jwt.Parse(str, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("token signature invalid")
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if i := claims["sid"].(float64); ok {
				idx := int(i)
				return t.keys[idx], nil
			}
		}
		return nil, errors.New("token signing invalid")
	})
	if err != nil || token.Claims.Valid() != nil || !token.Valid {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if v, ok := claims["uid"].(string); ok {
			return v, nil
		}
	}
	return "", errors.New("token invalid")
}
