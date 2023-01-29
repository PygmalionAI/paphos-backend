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

// Registers a new User.
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

	return c.Render(http.StatusCreated, r.JSON(user))
}

// Logs a user in.
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
