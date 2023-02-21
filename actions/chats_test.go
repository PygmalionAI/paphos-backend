package actions

import (
	"net/http"
	"paphos/models"

	"github.com/gobuffalo/nulls"
	"github.com/gofrs/uuid"
)

// Test_ChatsResource_List_WhenLoggedOut asserts that the List action is
// gated behind authentication.
func (as *ActionSuite) Test_ChatsResource_List_WhenLoggedOut() {
	res := as.JSON("/api/v1/chats").Get()

	as.Equal(http.StatusUnauthorized, res.Code)
}

// Test_ChatsResource_List asserts that the List action shows the Chats created
// by the User making the request, but not other Users.
func (as *ActionSuite) Test_ChatsResource_List() {
	as.LoadFixture("default")
	res := as.JSONAsUser("user@example.com", "/api/v1/chats").Get()

	// Should return HTTP 200 OK when properly authenticated.
	as.Equal(http.StatusOK, res.Code)

	// List action should not return Chats created by other Users.
	body := res.Body.String()
	as.NotContains(body, "Chat by Normal User #2")

	// List action should return Chats created by the current user.
	as.Contains(body, "Chat by Normal User")
}

// Test_ChatsResource_Create_WhenLoggedOut asserts that the Create action
// is gated behind authentication.
func (as *ActionSuite) Test_ChatsResource_Create_WhenLoggedOut() {
	params := ChatCreateParams{}
	res := as.JSON("/api/v1/chats").Post(params)

	as.Equal(http.StatusUnauthorized, res.Code)
}

// Test_ChatsResource_Create_WithBadData asserts that the Create action fails
// when called with bad data (count of character IDs != 1, or bad character ID).
func (as *ActionSuite) Test_ChatsResource_Create_WithBadData() {
	as.LoadFixture("default")

	params := ChatCreateParams{}
	params2 := ChatCreateParams{
		CharacterIDs: []uuid.UUID{
			uuid.Must(uuid.FromString("1df30357-d479-4861-893b-411fb7586dd1"))},
	}

	res := as.JSONAsUser("user@example.com", "/api/v1/chats").Post(params)
	res2 := as.JSONAsUser("user@example.com", "/api/v1/chats").Post(params2)

	// body := res.Body.String()

	// Can't assert on this since we're returning as a 500 (which gets masked when
	// running in prod/test mode).
	// as.Contains(body, "chats must have exactly one character for now")
	as.Equal(http.StatusInternalServerError, res.Code)
	as.Equal(http.StatusNotFound, res2.Code)
}

// Test_ChatsResource_Create asserts that the Create action successfully
// creates a Chat and the relevant ChatParticipants when given a single
// Character ID.
func (as *ActionSuite) Test_ChatsResource_Create() {
	as.LoadFixture("default")

	participant := &models.ChatParticipant{}
	originalCount, err := as.DB.Count(participant)
	as.NoError(err)

	character := &models.Character{}
	err = as.DB.Select("id").Where("name = ?", "Unlisted Character").First(character)
	as.NoError(err)

	params := ChatCreateParams{
		CharacterIDs: []uuid.UUID{character.ID},
		Name:         nulls.NewString("test chat"),
	}
	res := as.JSONAsUser("user@example.com", "/api/v1/chats").Post(params)

	// Assert that a HTTP 201 Created was correctly returned.
	as.Equal(http.StatusCreated, res.Code)

	newCount, err := as.DB.Count(participant)
	as.NoError(err)

	// Assert that there's two new ChatParticipants (one is the current user, the
	// other is the character passed in `params`).
	as.Equal(2, newCount-originalCount)
}

// Test_ChatsResource_Update_WhenLoggedOut asserts that the Update action
// is gated behind authentication.
func (as *ActionSuite) Test_ChatsResource_Update_WhenLoggedOut() {
	params := ChatUpdateParams{}
	res := as.JSON("/api/v1/chats/xxx").Put(params)

	as.Equal(http.StatusUnauthorized, res.Code)
}

// Test_ChatsResource_Update_WhenNoPermission asserts that the Update action
// does not allow updating a chat you didn't create.
func (as *ActionSuite) Test_ChatsResource_Update_WhenNoPermission() {
	as.LoadFixture("default")

	chat := &models.Chat{}
	as.DB.First(chat)

	params := ChatUpdateParams{Name: nulls.NewString("new chat name")}
	res := as.JSONAsUser("user2@example.com", "/api/v1/chats/"+chat.ID.String()).Put(params)

	as.Equal(http.StatusNotFound, res.Code)
}

// Test_ChatsResource_Update asserts that a User can update the name of a Chat
// he created.
func (as *ActionSuite) Test_ChatsResource_Update() {
	as.LoadFixture("default")

	chat := &models.Chat{}
	as.DB.First(chat)

	params := ChatUpdateParams{Name: nulls.NewString("new chat name")}
	res := as.JSONAsUser("user@example.com", "/api/v1/chats/"+chat.ID.String()).Put(params)
	as.Equal(http.StatusOK, res.Code)

	as.DB.Find(chat, chat.ID)
	as.Equal("new chat name", chat.Name.String)
}
