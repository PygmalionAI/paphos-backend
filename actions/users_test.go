package actions

import (
	"net/http"
	"paphos/models"
)

func (as *ActionSuite) Test_UsersLoginPost_WithBadCredentials() {
	as.LoadFixture("default")
	user := &models.User{Email: "user@example.com", Password: "wrong-password"}
	res := as.JSON("/api/v1/users/login").Post(user)

	// Attempting to log in with incorrect credentials should return a
	// HTTP Unauthorized response.
	as.Equal(http.StatusUnauthorized, res.Code)
}

func (as *ActionSuite) Test_UsersLoginPost_WithCorrectCredentials() {
	as.LoadFixture("default")
	// FIXME(11b): bcrypt returning mismatch between fixture hash and password so
	// this test is failing
	user := &models.User{Email: "user@example.com", Password: "123"}
	res := as.JSON("/api/v1/users/login").Post(user)

	// Attempting to log in with correct credentials should return a HTTP Created
	// response and a session JWT.
	as.Equal(http.StatusCreated, res.Code)
	as.Contains(res.Body.String(), "jwt")
}

func (as *ActionSuite) Test_UsersRegisterPost_EmailAlreadyInUse() {
	as.LoadFixture("default")
	user := &models.User{
		Email:                "user@example.com",
		DisplayName:          "Mr. New User",
		Password:             "123",
		PasswordConfirmation: "123",
	}
	res := as.JSON("/api/v1/users/register").Post(user)

	// Attempting to register with an email that's already in use should return
	// an error.
	as.Equal(http.StatusUnprocessableEntity, res.Code)
	as.Contains(res.Body.String(), "An account with this email already exists.")
}

func (as *ActionSuite) Test_UsersRegisterPost_MismatchedPasswords() {
	as.LoadFixture("default")
	user := &models.User{
		Email:                "newuser@example.com",
		DisplayName:          "Mr. New User",
		Password:             "123",
		PasswordConfirmation: "124",
	}
	res := as.JSON("/api/v1/users/register").Post(user)

	// Attempting to register with mismatched passwords should return an error.
	as.Equal(http.StatusUnprocessableEntity, res.Code)
	as.Contains(res.Body.String(), "Passwords do not match.")
}

func (as *ActionSuite) Test_UsersRegisterPost_OK() {
	as.LoadFixture("default")
	user := &models.User{
		Email:                "newuser@example.com",
		DisplayName:          "Mr. New User",
		Password:             "123",
		PasswordConfirmation: "123",
	}
	res := as.JSON("/api/v1/users/register").Post(user)

	// Attempting to register with valid data should return a HTTP Created
	// response.
	as.Equal(http.StatusCreated, res.Code)

	// And response should not contain sensitive fields.
	as.NotContains(res.Body.String(), "password")
	as.NotContains(res.Body.String(), "verification_token")
}
