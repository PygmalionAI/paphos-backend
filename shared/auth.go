package shared

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"paphos/models"
	"strings"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
)

var JWT_SECRET = []byte(os.Getenv("JWT_SECRET"))

// ExtractDataFromJWTMiddleware parses the JWT present in the Authorization
// header and injects the claims into the Buffalo session under the `claims` key
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

// ExtractUserUUIDFromContext fetches the UUID of the logged in user from the
// given `buffalo.Context`
func ExtractUserUUIDFromContext(c buffalo.Context) (uuid.UUID, error) {
	claims := c.Value("claims").(jwt.MapClaims)
	userFromClaims := claims["user"].(map[string]interface{})
	userUuidString := userFromClaims["id"].(string)
	return uuid.FromString(userUuidString)
}

// CreateSignedJWTStringForUser creates a signed session JWT for the given user.
func CreateSignedJWTStringForUser(u *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": models.UserFromJWT{
			ID:          u.ID,
			Email:       u.Email,
			DisplayName: u.DisplayName,
		},
		"iat": jwt.NewNumericDate(time.Now()),
		"exp": jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	})

	return token.SignedString(JWT_SECRET)
}
