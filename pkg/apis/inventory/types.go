package inventory

//go:generate counterfeiter -o ./mocks/identifier.go . Identifier

// Identifier of an inventory
type Identifier interface {
	GUID() string
}

type identifier struct {
	guid string
}

func (i *identifier) GUID() string {
	return i.guid
}

// Provider of inventory information
type Provider interface {
	Data() *data
}

// Incomplete inventory object
type Incomplete interface {
	Provider
}

// Complete inventory object
type Complete interface {
	Identifier
	Incomplete
}

// Item of an inventory
type Item struct {
	ID     string  `json:"id"`
	Amount uint64  `json:"amount"`
	Subset float64 `json:"subset"`
}

type data struct {
	Items []*Item `json:"items"`
}

// NewIncomplete returns a new incomplete inventory object instance
func NewIncomplete(items []*Item) Incomplete {
	return &data{items}
}

func (d *data) Data() *data {
	return d
}
