package token

import "context"

// EventPayload of a token in the event system of the API.
type EventPayload struct {
	*Info
}

// EventPayloadFromContext returns a necessary event payload of token
// information from a given context.
func EventPayloadFromContext(ctx context.Context) (*EventPayload, error) {
	tok, err := FromContext(ctx)
	if err != nil {
		return nil, err
	}

	return &EventPayload{
		&Info{
			User: tok.Data().User,
		},
	}, nil
}
