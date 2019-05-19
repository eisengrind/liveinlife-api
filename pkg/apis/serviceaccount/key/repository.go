package key

//go:generate counterfeiter -o ./mocks/repository.go . Repository

import "context"

// Repository represents methods to save and read data regarding keys
type Repository interface {
	Get(context.Context, Identifier) (Complete, error)
	Update(context.Context, Complete) error
	Create(context.Context, Incomplete) (Complete, error)
	Delete(context.Context, Identifier) error
}
