package inventory

//go:generate counterfeiter -o ./mocks/repository.go . Repository

import "context"

// Repository to store inventory information
type Repository interface {
	Get(context.Context, Identifier) (Complete, error)
	Create(context.Context) (Complete, error)
	AddItem(context.Context, Identifier, *Item) error
	RemoveItem(context.Context, Identifier, *Item) error
	Delete(context.Context, Identifier) error
}
