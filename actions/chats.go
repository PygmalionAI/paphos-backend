package actions

import (
	"fmt"
	"net/http"
	"paphos/models"
	"paphos/shared"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
)

type ChatsResource struct {
	buffalo.Resource
}

// ChatCreateParams defines the request body required when creating a Chat.
type ChatCreateParams struct {
	// NOTE(11b): At some point, I'd like for multi-user chats to be a feature.
	// That's a lot of complexity though (will need chat admins, a way to stop
	// trolls from adding you to chats you don't want, leaving a chat, etc.) so
	// we'll leave that for later.
	//
	// UserIDs      []uuid.UUID  `json:"user_ids" example:"[\"e3609a19-18d9-4240-b300-d92291a48bc5\"]"`

	CharacterIDs []uuid.UUID  `json:"character_ids" example:"[\"1379859c-1113-48b4-97e4-36e5fe375e12\"]"`
	Name         nulls.String `json:"name" example:"A test chat"`
}

// Create adds a Chat to the DB alongside the relevant ChatParticipants.
func (v ChatsResource) Create(c buffalo.Context) error {
	// Parse request body into the `params` struct.
	params := &ChatCreateParams{}
	if err := c.Bind(params); err != nil {
		return err
	}

	// TODO(11b): Make this a 422. Not sure what's the cleanest way to do that at
	// the action-level since we're not running a model validator here.
	if len(params.CharacterIDs) != 1 {
		return fmt.Errorf("chats must have exactly one character for now")
	}

	// Get the database transaction from the request context.
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Create Chat from request params and request context information.
	userUUID, err := shared.ExtractUserUUIDFromContext(c)
	if err != nil {
		return err
	}

	chatUUID, err := uuid.NewV4()
	if err != nil {
		return err
	}

	participants := models.ChatParticipants{}
	chat := &models.Chat{
		ID:      chatUUID,
		OwnerID: userUUID,
		Name:    params.Name,
	}

	{
		verrs, err := tx.ValidateAndCreate(chat)
		if err != nil {
			return err
		}
		if verrs.HasAny() {
			return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
		}
	}

	// Now, create the ChatParticipant(s). First the User making the request, then
	// the Characters passed in via params.
	{
		participant := &models.ChatParticipant{
			UserID: nulls.NewUUID(userUUID),
			ChatID: chatUUID,
		}

		verrs, err := tx.ValidateAndCreate(participant)
		if err != nil {
			return err
		}
		if verrs.HasAny() {
			return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
		}

		participants = append(participants, *participant)
	}

	for _, characterID := range params.CharacterIDs {
		// First, make sure this Character exists and is accessible by this User.
		character := &models.Character{}
		query := tx.
			Where("id = ?", characterID).
			Scope(models.CharactersAccessibleByUser(userUUID))

		exists, err := query.Exists(character)
		if err != nil {
			return err
		}
		if !exists {
			return c.Error(http.StatusNotFound, fmt.Errorf("character not found"))
		}

		// If everything is OK, add as a participant to the Chat.
		participant := &models.ChatParticipant{
			CharacterID: nulls.NewUUID(characterID),
			ChatID:      chatUUID,
		}

		verrs, err := tx.ValidateAndCreate(participant)
		if err != nil {
			return err
		}
		if verrs.HasAny() {
			return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
		}

		participants = append(participants, *participant)
	}

	chat.Participants = participants
	return c.Render(http.StatusCreated, r.JSON(chat))
}

// List returns all the Chats that the User who is making the HTTP request has
// created.
func (v ChatsResource) List(c buffalo.Context) error {
	// Get the database transaction from the request context.
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	userUUID, err := shared.ExtractUserUUIDFromContext(c)
	if err != nil {
		return err
	}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	query := tx.PaginateFromParams(c.Params())

	chats := &models.Chats{}
	if err := query.
		EagerPreload().
		Scope(models.ChatsAccessibleByUser(userUUID)).
		All(chats); err != nil {
		return err
	}

	return c.Render(200, r.JSON(chats))
}

// ChatUpdateParams defines the request body required when updating a Chat.
type ChatUpdateParams struct {
	Name nulls.String `json:"name" example:"Renamed chat"`
}

// Update allows updating the name of a given Chat.
func (v ChatsResource) Update(c buffalo.Context) error {
	// Parse request body into the `params` struct.
	params := &ChatUpdateParams{}
	if err := c.Bind(params); err != nil {
		return err
	}

	// Get the database transaction from the request context.
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Find the Chat referenced by the ID given in the path parameters, making
	// sure that it is accessible to the current user.
	userUUID, err := shared.ExtractUserUUIDFromContext(c)
	if err != nil {
		return err
	}

	chat := &models.Chat{}
	query := tx.
		// Select("id", "owner_id").
		Scope(models.ChatsAccessibleByUser(userUUID)).
		Where("id = ?", c.Param("chat_id"))
	if err := query.First(chat); err != nil {
		return err
	}
	if chat.ID.IsNil() {
		return c.Error(http.StatusNotFound, fmt.Errorf("chat not found"))
	}

	// Update it with the incoming params.
	chat.Name = params.Name

	verrs, err := tx.ValidateAndUpdate(chat, "id", "owner_id")
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
	}

	return c.Render(http.StatusOK, r.JSON(chat))
}
