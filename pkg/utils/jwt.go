package utils

import "github.com/golang-jwt/jwt/v5"

type Customclaims struct {
	UserName string `json:"username"`
	jwt.RegisteredClaims
}

type Setting struct {
	Secret []byte
	Customclaims
}

func GenerateToken(setting *Setting) (string, error) {
	claims := setting.Customclaims

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(setting.Secret)

	return tokenStr, err
}

func ParseToken(token string, secret []byte) (*Customclaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Customclaims{}, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Customclaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
