package user

import (
	"github.com/51st-state/api/pkg/event"
)

// CreatedEventID of an user object
const CreatedEventID event.ID = "user_created"

// CreatedEvent of an user object
type CreatedEvent struct {
	Meta *event.PayloadMeta `json:"meta"`
	Data Complete           `json:"data"`
}

// UpdatedEventID of an user object
const UpdatedEventID event.ID = "user_updated"

// UpdatedEvent of an user object
type UpdatedEvent struct {
	Meta *event.PayloadMeta `json:"meta"`
	Data Complete           `json:"data"`
}

// DeletedEventID of an user object
const DeletedEventID event.ID = "user_deleted"

// DeletedEvent of an user object
type DeletedEvent struct {
	Meta *event.PayloadMeta `json:"meta"`
	Data Identifier         `json:"data"`
}

// PasswordSetEventID of an user object
const PasswordSetEventID event.ID = "user_password_set"

// PasswordSetEvent of an user object
type PasswordSetEvent struct {
	Meta *event.PayloadMeta `json:"meta"`
	Data CompletePassword   `json:"data"`
}
