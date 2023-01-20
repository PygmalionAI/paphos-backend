package shared

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/golang-jwt/jwt/v4"
)

var JWT_SECRET = []byte(os.Getenv("JWT_SECRET"))

func ExtractDataFromJWTMiddleware(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 {
			return c.Error(http.StatusUnauthorized, errors.New("Unauthorized"))
		}

		token, err := jwt.Parse(authHeaderParts[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return JWT_SECRET, nil
		})

		if err != nil {
			return c.Error(http.StatusUnauthorized, errors.New("Unauthorized"))
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("claims", claims)
			return next(c)
		}

		return c.Error(http.StatusUnauthorized, errors.New("Unauthorized"))
	}
}
