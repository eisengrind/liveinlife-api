package email

import (
	"time"

	"github.com/51st-state/api/pkg/event"
)

// NewEvent of an email.
// Fired from producers who want to send an email to receivers.
const NewEvent event.ID = "email_new"

// CreatedEvent of an email.
// Fired when an email was sent successfully.
const CreatedEvent event.ID = "email_created"

// Email to send to receivers
type Email struct {
	Header *Header `json:"header"`
	Body   string  `json:"body"`
}

// Header of an Email
type Header struct {
	From      string      `json:"from"`
	To        []string    `json:"to"`
	Cc        []Recipient `json:"cc"`
	Bcc       []Recipient `json:"bcc"`
	Subject   string      `json:"subject"`
	CreatedAt time.Time   `json:"created_at"`
}

// Recipient of emails who are sent to third party people (Cc, Bcc)
type Recipient struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// New email for the event queue
func New(from, subject, body string, to []string, cc, bcc []Recipient) *Email {
	return &Email{
		&Header{
			from,
			to,
			cc,
			bcc,
			subject,
			time.Now(),
		},
		body,
	}
}
