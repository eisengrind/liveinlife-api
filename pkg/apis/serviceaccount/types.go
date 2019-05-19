package serviceaccount

//go:generate counterfeiter -o ./mocks/identifier.go . Identifier

// Identifier of a service account
type Identifier interface {
	GUID() string
}

type identifier struct {
	guid string
}

// NewIdentifier creates a new identifier object
func NewIdentifier(guid string) Identifier {
	return &identifier{guid}
}

func (i *identifier) GUID() string {
	return i.guid
}

// Provider provides methods for the incomplete service account object
type Provider interface {
	Data() *data
}

// Incomplete represents an incomplete service account object
type Incomplete interface {
	Provider
}

// Complete represens a complete service account object
type Complete interface {
	Identifier
	Incomplete
}

type complete struct {
	Identifier
	Incomplete
}

type data struct {
	Name        string
	Description string
}

// NewIncomplete creates a new incomplete service account object
func NewIncomplete(name, description string) Incomplete {
	return &data{
		name,
		description,
	}
}

func (d *data) Data() *data {
	return d
}

func (d *data) SetName(to string) *data {
	d.Name = to
	return d
}

func (d *data) SetDescription(to string) *data {
	d.Description = to
	return d
}
