package actions

import (
	"fmt"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"paphos/models"
	"paphos/shared"
)

// Commented out for now to avoid data leakage. Adjust when we actually need
// this data in the front-end.
/* func UsersShowGet(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	user := &models.User{}
	if err := tx.Find(user, c.Param("user_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	return c.Render(200, r.JSON(user))
} */

type RegisterParameters struct {
	Email                string `json:email example: example@mail.com`
	DisplayName          string `json:display_name example: My Name`
	Password             string `json:password example: pass123`
	PasswordConfirmation string `json:password_confirmation example: pass123`
}

type LoginParameters struct {
	Email    string `json:email example: example@mail.com`
	Password string `json:password example: pass123`
}

// UsersRegisterPost godoc
// POST
// @Summary      Registers a new User
// @Description  Registers a new User in the DB.
// @Tags         Users
// @Param        register  body  RegisterParameters  true  "Register Body"
// @Produce      json
// @Success      200
// @Router       /users/register [POST]
func UsersRegisterPost(c buffalo.Context) error {
	// Allocate an empty User
	user := &models.User{}

	// Bind user to the request body payload
	if err := c.Bind(user); err != nil {
		return err
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Validate the data from the request
	verrs, err := user.Create(tx)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
	}

	return c.Render(http.StatusCreated, r.JSON(struct {
		ID          uuid.UUID `json:"id"`
		Email       string    `json:"email"`
		DisplayName string    `json:"display_name"`
	}{user.ID, user.Email, user.DisplayName}))
}

// UsersLoginPost godoc
// POST
// @Summary      Log the user in.
// @Description  Logs a user in by returning a session JWT.
// @Tags         Users
// @Param        login  body  LoginParameters  true  "Login Body"
// @Accept       json
// @Produce      json
// @Success      200
// @Router       /users/login [POST]
func UsersLoginPost(c buffalo.Context) error {
	user := &models.User{}

	if err := c.Bind(user); err != nil {
		return err
	}

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	populated_user, err := user.Authorize(tx)
	if err != nil {
		verrs := validate.NewErrors()
		verrs.Add("root", "Invalid email or password.")
		return c.Error(http.StatusUnauthorized, verrs)
	}

	tokenString, err := shared.CreateSignedJWTStringForUser(populated_user)
	if err != nil {
		return err
	}

	return c.Render(http.StatusCreated, r.JSON(struct {
		ID          uuid.UUID `json:"id"`
		Email       string    `json:"email"`
		DisplayName string    `json:"display_name"`
		JWT         string    `json:"jwt"`
	}{populated_user.ID, populated_user.Email, populated_user.DisplayName, tokenString}))
}
