package models

import (
	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

//
// Database model
//

// ChatParticipant is used by pop to map the chat_participants database table to
// Go code.
type ChatParticipant struct {
	// Unfortunately, it looks like pop requires us to add this useless ID column.
	// https://github.com/gobuffalo/pop/blob/ec9229dbf7d7ccd49926ead635685e27b6fa0fb4/model.go#L41
	ID uuid.UUID `json:"-" db:"id"`

	UserID      nulls.UUID `json:"user_id" db:"user_id"`
	CharacterID nulls.UUID `json:"character_id" db:"character_id"`
	ChatID      uuid.UUID  `json:"-" db:"chat_id"`
	Chat        Chat       `json:"-" belongs_to:"chat"`
}

type ChatParticipants []ChatParticipant

//
// Custom validations
//

type onlyOneIDPresent struct {
	Name        string
	UserID      nulls.UUID
	CharacterID nulls.UUID
}

// IsValid performs the validation check to make sure exactly one ID (User or
// Character) has been given.
func (v *onlyOneIDPresent) IsValid(errors *validate.Errors) {
	if v.UserID.Valid && v.CharacterID.Valid {
		errors.Add(validators.GenerateKey(v.Name),
			"Participant must be either a User or a Character, not both.")
		return
	}

	if v.UserID.Valid || v.CharacterID.Valid {
		return
	}

	errors.Add(validators.GenerateKey(v.Name),
		"Participant must be either a User or a Character, but no IDs were given.")
}

//
// Validations
//

// Validate gets run every time a "pop.Validate*" method (pop.ValidateAndSave,
// pop.ValidateAndCreate, pop.ValidateAndUpdate) gets called.
func (cp *ChatParticipant) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&onlyOneIDPresent{
			Name:        "UserID",
			UserID:      cp.UserID,
			CharacterID: cp.CharacterID,
		},
		&validators.UUIDIsPresent{Field: cp.ChatID, Name: "RoomID"},
	), nil
}
