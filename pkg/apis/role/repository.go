package role

import "context"

// Repository for storage of role information
//go:generate counterfeiter -o ./mocks/repository.go . Repository
type Repository interface {
	Get(context.Context, Identifier) (Complete, error)
	Update(context.Context, Complete) error
	Create(context.Context, Complete) error
	Delete(context.Context, Identifier) error
}
