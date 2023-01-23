package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// Character is used by pop to map your characters database table to your go code.
type Character struct {
	ID uuid.UUID `json:"id" db:"id"`

	Name          string       `json:"name" db:"name"`
	Description   string       `json:"description" db:"description"`
	AvatarID      nulls.String `json:"avatar_id" db:"avatar_id"`
	Greeting      string       `json:"greeting" db:"greeting"`
	Persona       string       `json:"persona" db:"persona"`
	WorldScenario nulls.String `json:"world_scenario" db:"world_scenario"`
	ExampleChats  nulls.String `json:"example_chats" db:"example_chats"`
	Visibility    string       `json:"visibility" db:"visibility"`

	Creator   User      `json:"-" belongs_to:"user"`
	CreatorID uuid.UUID `json:"-" db:"creator_id"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (c Character) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// Characters is not required by pop and may be deleted
type Characters []Character

// String is not required by pop and may be deleted
func (c Characters) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// CharactersVisibleToUser defines a scope to return only Characters that are visible to
// the user with the given UUID.
func CharactersVisibleToUser(userUuid uuid.UUID) pop.ScopeFunc {
	return func(q *pop.Query) *pop.Query {
		q = q.Where("visibility = 'public' OR creator_id = ?", userUuid)
		return q
	}
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (c *Character) Validate(tx *pop.Connection) (*validate.Errors, error) {
	// TODO(11b): All these hardcoded values might be better off at the top of the
	// file, or being defined somewhere else entirely.
	VALID_VISIBILITY_VALUES := []string{"public", "unlisted", "private"}

	return validate.Validate(
		&validators.StringIsPresent{Field: c.Name, Name: "Name"},
		&validators.StringLengthInRange{Field: c.Name, Name: "Name", Min: 1, Max: 32},

		&validators.StringIsPresent{Field: c.Description, Name: "Description"},
		&validators.StringLengthInRange{Field: c.Description, Name: "Description", Min: 12, Max: 64},

		&validators.StringIsPresent{Field: c.Greeting, Name: "Greeting"},
		&validators.StringLengthInRange{Field: c.Greeting, Name: "Greeting", Min: 2, Max: 1024},

		&validators.StringIsPresent{Field: c.Persona, Name: "Persona"},
		&validators.StringLengthInRange{Field: c.Persona, Name: "Persona", Min: 12, Max: 1024},

		&validators.StringLengthInRange{Field: c.WorldScenario.String, Name: "Scenario", Min: 0, Max: 1024},
		&validators.StringLengthInRange{Field: c.ExampleChats.String, Name: "Example chats", Min: 0, Max: 1024},

		&validators.StringInclusion{Field: c.Visibility, Name: "Visibility", List: VALID_VISIBILITY_VALUES},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (c *Character) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (c *Character) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
