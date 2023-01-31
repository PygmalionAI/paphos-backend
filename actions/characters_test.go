package actions

import (
	"net/http"
	"paphos/models"
	"paphos/shared"

	"github.com/gobuffalo/httptest"
)

// JSONAsUser creates a test *httptest.JSON request while authenticated as the
// User identified by the given `userEmail`.
func (as *ActionSuite) JSONAsUser(userEmail string, u string, args ...interface{}) *httptest.JSON {
	user := &models.User{}
	err := as.DB.Where("email = ?", userEmail).First(user)
	as.NoError(err)

	jwt, err := shared.CreateSignedJWTStringForUser(user)
	as.NoError(err)
	req := as.JSON(u, args...)
	req.Headers["Authorization"] = "Bearer " + jwt

	return req
}

// TODO(11b): ^ move out of here, we'll likely use this elsewhere as well

// Test_CharactersResource_List_WhenLoggedOut asserts that the List action is
// gated behind authentication.
func (as *ActionSuite) Test_CharactersResource_List_WhenLoggedOut() {
	as.LoadFixture("default")
	res := as.JSON("/api/v1/characters").Get()

	as.Equal(http.StatusUnauthorized, res.Code)
}

// Test_CharactersResource_List asserts that the List action does not serialize
// useless fields, and does not include private/unlisted Characters.
func (as *ActionSuite) Test_CharactersResource_List() {
	as.LoadFixture("default")
	res := as.JSONAsUser("user@example.com", "/api/v1/characters").Get()

	// Should return HTTP 200 OK when properly authenticated.
	as.Equal(http.StatusOK, res.Code)

	// List action should only return the minimal set of required fields.
	body := res.Body.String()
	as.NotContains(body, "Hello! I'm a public character.")                      // greeting
	as.NotContains(body, "An example character marked with public visibility.") // persona

	// A note for the above assertions: ideally we'd completely drop the fields
	// from the response JSON, then we could simply assert that the field names
	// don't exist and this test would be a lot more robust. However, as a Golang
	// newb I couldn't find a way to cleanly do that without resorting to either
	// lots of code duplication or invoking reflection. My approach was then to
	// just drop the contents of the fields instead and keep using the same model
	// structs as a middle-of-the-road trade-off.

	// List action should not return the private/unlisted Characters of other
	// Users.
	as.NotContains(body, "Private Character")
}

// Test_CharactersResource_Show_WhenLoggedOut asserts that the Show action is
// gated behind authentication.
func (as *ActionSuite) Test_CharactersResource_Show_WhenLoggedOut() {
	as.LoadFixture("default")

	character := &models.Character{}
	err := as.DB.Select("id").First(character)
	as.NoError(err)
	res := as.JSON("/api/v1/characters/" + character.ID.String()).Get()

	as.Equal(http.StatusUnauthorized, res.Code)
}

// Test_CharactersResource_Show_WhenPrivate asserts that ...
func (as *ActionSuite) Test_CharactersResource_Show_WhenPrivate() {
	as.LoadFixture("default")

	character := &models.Character{}
	err := as.DB.Select("id").Where("visibility = 'private'").First(character)
	as.NoError(err)
	endpoint := "/api/v1/characters/" + character.ID.String()

	res := as.JSONAsUser("user@example.com", endpoint).Get()
	otherRes := as.JSONAsUser("user2@example.com", endpoint).Get()

	// Only creators should be able to access their private characters.
	as.Equal(http.StatusNotFound, res.Code)
	as.Equal(http.StatusOK, otherRes.Code)
}

// Test_CharactersResource_Create_WhenLoggedOut asserts that the Create action is
// gated behind authentication.
func (as *ActionSuite) Test_CharactersResource_Create_WhenLoggedOut() {
	as.LoadFixture("default")

	character := models.Character{}
	res := as.JSON("/api/v1/characters").Post(character)

	as.Equal(http.StatusUnauthorized, res.Code)
}

// Test_CharactersResource_Create asserts that the Create action is successfully
// creates a Character.
func (as *ActionSuite) Test_CharactersResource_Create() {
	as.LoadFixture("default")

	character := models.Character{
		Name:        "name",
		Description: "example description",
		Persona:     "example persona",
		Greeting:    "greeting",
		Visibility:  "public",
	}
	res := as.JSONAsUser("user@example.com", "/api/v1/characters").Post(character)

	as.Equal(http.StatusCreated, res.Code)

	// TODO(11b): Consider asserting that the proper `creator_id` was added to the
	// database.
	/*
		err := json.Unmarshal(res.Body.Bytes(), &character)
		as.NoError(err)

		as.DB.Select("creator_id").Where("id = ?", character.ID).First(character)
		as.Equal(...)
	*/
}

// Test_CharactersResource_Update_WhenLoggedOut asserts that the Update action is
// gated behind authentication.
func (as *ActionSuite) Test_CharactersResource_Update_WhenLoggedOut() {
	as.LoadFixture("default")

	character := models.Character{}
	res := as.JSON("/api/v1/characters/xxx").Put(character)

	as.Equal(http.StatusUnauthorized, res.Code)
}

// Test_CharactersResource_Update_WhenNotCreator asserts that only a Character's
// creator can edit it.
func (as *ActionSuite) Test_CharactersResource_Update_WhenNotCreator() {
	as.LoadFixture("default")

	originalCharacter := &models.Character{}
	err := as.DB.Where("name = 'Public Character'").First(originalCharacter)
	as.NoError(err)

	character := models.Character{}
	res := as.JSONAsUser("user2@example.com", "/api/v1/characters/"+originalCharacter.ID.String()).Put(character)

	as.Equal(http.StatusForbidden, res.Code)
}

// Test_CharactersResource_Update asserts that a Character can be successfully
// edited by its creator.
func (as *ActionSuite) Test_CharactersResource_Update() {
	as.LoadFixture("default")

	// Fetch the original Character from the DB.
	dbCharacter := &models.Character{}
	err := as.DB.Where("name = 'Public Character'").First(dbCharacter)
	as.NoError(err)

	// Update its name.
	reqCharacter := struct {
		Name string `json:"name"`
	}{Name: "Updated Character"}
	res := as.JSONAsUser("user@example.com", "/api/v1/characters/"+dbCharacter.ID.String()).Put(reqCharacter)

	as.Equal(http.StatusOK, res.Code)

	// Finally, assert that the name in the DB matches the name we just updated with.
	as.DB.Select("name").Where("id = ?", dbCharacter.ID.String()).First(dbCharacter)
	as.Equal(reqCharacter.Name, dbCharacter.Name)
}
