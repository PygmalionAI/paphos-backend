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

// User is used by pop to map your users database table to your go code.
type User struct {
	ID uuid.UUID `json:"id" db:"id"`

	Email          string `json:"email" db:"email"`
	HashedPassword string `json:"hashed_password" db:"hashed_password"`
	DisplayName    string `json:"display_name" db:"display_name"`
	Role           string `json:"-" db:"role"`

	VerificationToken  nulls.String `json:"-" db:"verification_token"`
	PasswordResetToken nulls.String `json:"-" db:"password_reset_token"`

	// Used during registration.
	Password             string `json:"-" db:"-"`
	PasswordConfirmation string `json:"-" db:"-"`

	LastLogin nulls.Time `json:"last_login" db:"last_login"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (u User) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Users is not required by pop and may be deleted
type Users []User

// String is not required by pop and may be deleted
func (u Users) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (u *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
	// TODO(11b): Move elsewhere?
	VALID_ROLE_VALUES := []string{"admin", "user"}

	return validate.Validate(
		&validators.EmailIsPresent{Field: u.Email, Name: "Email"},
		&validators.StringLengthInRange{Field: u.Email, Name: "Email", Min: 8, Max: 96},
		&validators.StringIsPresent{Field: u.HashedPassword, Name: "HashedPassword"},

		&validators.StringIsPresent{Field: u.DisplayName, Name: "DisplayName"},
		&validators.StringLengthInRange{Field: u.DisplayName, Name: "DisplayName", Min: 12, Max: 64},

		&validators.StringInclusion{Field: u.Role, Name: "Role", List: VALID_ROLE_VALUES},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (u *User) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (u *User) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
