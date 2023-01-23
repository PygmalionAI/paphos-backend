package actions

import (
	"fmt"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"

	"paphos/models"
	"paphos/shared"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (Character)
// DB Table: Plural (characters)
// Resource: Plural (Characters)
// Path: Plural (/characters)
// View Template Folder: Plural (/templates/characters/)

// CharactersResource is the resource for the Character model
type CharactersResource struct {
	buffalo.Resource
}

// List gets all Characters. This function is mapped to the path
// GET /characters
func (v CharactersResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	characters := &models.Characters{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve visible Characters from the DB
	userUuid, err := shared.ExtractUserUUIDFromContext(c)
	if err != nil {
		return err
	}

	if err := q.Scope(models.CharactersVisibleToUser(userUuid)).All(characters); err != nil {
		return err
	}

	return c.Render(200, r.JSON(characters))
}

// Show gets the data for one Character. This function is mapped to
// the path GET /characters/{character_id}
func (v CharactersResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Character
	character := &models.Character{}

	// To find the Character the parameter character_id is used.
	userUuid, err := shared.ExtractUserUUIDFromContext(c)
	if err != nil {
		return err
	}

	// FIXME: leaking data
	// if err := tx.Scope(models.CharactersVisibleToUser(userUuid)).Find(character, c.Param("character_id")); err != nil {
	if err := tx.Scope(models.CharactersVisibleToUser(userUuid)).Where("id = ?", c.Param("character_id")).First(character); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	return c.Render(200, r.JSON(character))
}

// Create adds a Character to the DB. This function is mapped to the
// path POST /characters
func (v CharactersResource) Create(c buffalo.Context) error {
	// Allocate an empty Character
	character := &models.Character{}

	// Bind character to the html form elements
	if err := c.Bind(character); err != nil {
		return err
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Attach ID of the creator to the character
	userUuid, err := shared.ExtractUserUUIDFromContext(c)
	if err != nil {
		return err
	}

	character.CreatorID = userUuid

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(character)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
	}

	return c.Render(http.StatusCreated, r.JSON(character))
}

// Update changes a Character in the DB. This function is mapped to
// the path PUT /characters/{character_id}
func (v CharactersResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Character
	character := &models.Character{}

	if err := tx.Find(character, c.Param("character_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Bind Character to the html form elements
	if err := c.Bind(character); err != nil {
		return err
	}

	verrs, err := tx.ValidateAndUpdate(character)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
	}

	return c.Render(http.StatusOK, r.JSON(character))
}

// Destroy deletes a Character from the DB. This function is mapped
// to the path DELETE /characters/{character_id}
func (v CharactersResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Character
	character := &models.Character{}

	// To find the Character the parameter character_id is used.
	if err := tx.Find(character, c.Param("character_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := tx.Destroy(character); err != nil {
		return err
	}

	return c.Render(http.StatusOK, r.JSON(character))
}
