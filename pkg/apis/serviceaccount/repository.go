package serviceaccount

//go:generate counterfeiter -o ./mocks/repository.go . Repository

import "context"

// Repository to manage the storage of service accounts
type Repository interface {
	Get(context.Context, Identifier) (Complete, error)
	Create(context.Context, Incomplete) (Complete, error)
	Update(context.Context, Complete) error
	Delete(context.Context, Identifier) error
}
