package preselect

import "context"

// Repository for preselect objects
type Repository interface {
	GetNext(context.Context) (Complete, error)
	GetLeft(context.Context) (uint64, error)
	Create(context.Context, ...Complete) error
	Update(context.Context, ...Complete) error
}
