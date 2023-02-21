package models

import (
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/gofrs/uuid"
)

//
// Database model
//

// Chat is used by pop to map the chats database table to Go code.
type Chat struct {
	ID      uuid.UUID `json:"id" db:"id"`
	OwnerID uuid.UUID `json:"owner_id" db:"owner_id"`

	Name         nulls.String     `json:"name" db:"name"`
	Participants ChatParticipants `json:"participants" has_many:"chat_participants"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
