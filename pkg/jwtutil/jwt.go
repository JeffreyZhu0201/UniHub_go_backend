package jwtutil

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Claims holds JWT payload.
type Claims struct {
	UserID uint  `json:"user_id"`
	RoleID uint  `json:"role_id"`
	OrgID  *uint `json:"org_id"`
	jwt.RegisteredClaims
}

// Generate signs a token with given secret and ttl duration.
func Generate(secret string, ExpireHours int, userID, roleID uint) (string, error) {
	ttl := time.Duration(ExpireHours) * time.Hour
	claims := Claims{
		UserID: userID,
		RoleID: roleID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// Parse validates token and returns claims.
func Parse(tokenStr, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}
