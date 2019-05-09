package topgenerator

import "context"

// Repository for top storage
type Repository interface {
	Get(context.Context, Identifier) (Complete, error)
	Upsert(context.Context, Complete) error
}
