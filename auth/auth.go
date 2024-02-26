package auth

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/spf13/cast"
)

var JWT_KEY = []byte("super_secret")

func GenerateToken(userID string) string {
	token, err := jwt.NewBuilder().Issuer("github.com/iskaa02").IssuedAt(time.Now()).
		Claim("id", userID).
		Build()
	if err != nil {
		return ""
	}
	signed, err := jwt.Sign(token, jwt.WithKey(jwa.HS512, JWT_KEY))
	if err != nil {
		fmt.Printf("failed to sign token: %s\n", err)
		return ""
	}
	return string(signed)
}

func ParseToken(token string) (string, error) {
	verifiedToken, err := jwt.Parse([]byte(token), jwt.WithKey(jwa.HS512, JWT_KEY))
	if err != nil {
		fmt.Printf("failed to verify JWS: %s\n", err)
		return "", err
	}
	id, ok := verifiedToken.Get("id")
	if !ok {
		return "", fmt.Errorf("id not found")
	}
	return cast.ToString(id), nil
}

func isAdmin(token string, db *sqlx.DB) (string, bool) {
	id, err := ParseToken(token)
	if err != nil {
		return id, false
	}
	exists := false
	err = db.Get(&exists, "SELECT 1 FROM user_admin WHERE id=$1", id)
	if err != nil {
		return id, false
	}
	return id, exists
}

func isAuthenticated(token string) (string, bool) {
	id, _ := ParseToken(token)
	if id != "" {
		return id, true
	}
	return "", false
}
