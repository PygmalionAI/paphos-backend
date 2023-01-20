package models

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// User is used by pop to map your users database table to your go code.
type User struct {
	ID uuid.UUID `json:"id" db:"id"`

	Email          string `json:"email" db:"email"`
	HashedPassword string `json:"hashed_password" db:"hashed_password"`
	DisplayName    string `json:"display_name" db:"display_name"`
	Role           string `json:"role" db:"role"`

	VerificationToken  nulls.String `json:"-" db:"verification_token"`
	PasswordResetToken nulls.String `json:"-" db:"password_reset_token"`

	// Used during registration.
	Password             string `json:"password" db:"-"`
	PasswordConfirmation string `json:"password_confirmation" db:"-"`

	LastLogin nulls.Time `json:"last_login" db:"last_login"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}

// Marshals a User struct but only going over the public fields.
// https://stackoverflow.com/a/31374980
func (u *User) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID          uuid.UUID `json:"id"`
		Email       string    `json:"email"`
		DisplayName string    `json:"display_name"`
	}{u.ID, u.Email, u.DisplayName})
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

// Validates and creates a new User.
func (u *User) Create(tx *pop.Connection) (*validate.Errors, error) {
	u.Email = strings.ToLower(u.Email)
	u.Role = "user"

	pwdHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return validate.NewErrors(), errors.WithStack(err)
	}

	u.HashedPassword = string(pwdHash)
	return tx.ValidateAndCreate(u)
}

type EmailNotTaken struct {
	Name  string
	Field string
	tx    *pop.Connection
}

// Checks user's password for logging in.
func (u *User) Authorize(tx *pop.Connection) (error, *User) {
	err := tx.Select("id", "display_name", "hashed_password").Where("email = ?", strings.ToLower(u.Email)).First(u)
	if err != nil {
		return err, nil
	}

	// Confirm that the given password matches the hashed password from the DB.
	err = bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(u.Password))
	if err != nil {
		return err, nil
	}

	return nil, u
}

// IsValid performs the validation check for unique emails
func (v *EmailNotTaken) IsValid(errors *validate.Errors) {
	query := v.tx.Where("email = ?", v.Field)
	queryUser := User{}
	exists, err := query.Exists(&queryUser)
	if err != nil {
		errors.Add(validators.GenerateKey(v.Name), "We couldn't verify whether an account with this email already exists, please try again later.")
		return
	}

	if exists {
		errors.Add(validators.GenerateKey(v.Name), "An account with this email already exists.")
	}
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (u *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.EmailIsPresent{Field: u.Email, Name: "Email"},
		&validators.StringLengthInRange{Field: u.Email, Name: "Email", Min: 8, Max: 96},
		&EmailNotTaken{Name: "Email", Field: u.Email, tx: tx},

		&validators.StringIsPresent{Field: u.Password, Name: "Password"},
		&validators.StringIsPresent{Field: u.PasswordConfirmation, Name: "Password confirmation"},
		&validators.StringsMatch{
			Field:   u.Password,
			Field2:  u.PasswordConfirmation,
			Name:    "Password",
			Message: "Passwords do not match.",
		},

		&validators.StringIsPresent{Field: u.DisplayName, Name: "DisplayName"},
		&validators.StringLengthInRange{Field: u.DisplayName, Name: "DisplayName", Min: 2, Max: 64},
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
