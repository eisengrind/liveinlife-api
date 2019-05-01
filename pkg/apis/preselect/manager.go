package preselect

// Manager for pre-selection management
type Manager struct {
	repository Repository
}

// NewManager for pre-selection management
func NewManager(r Repository) *Manager {
	return &Manager{
		r,
	}
}

/*// GetLeft preselect objects
func (m *Manager) GetLeft(ctx context.Context) (uint64, error) {
	return m.repository.GetLeft(ctx)
}

// GetNext preselect object
func (m *Manager) GetNext(ctx context.Context) (Complete, error) {
	return m.repository.GetNext(ctx)
}

// Create pre-selections
func (m *Manager) Create(ctx context.Context, c ...Complete) error {
	return m.repository.Create(ctx, c...)
}

// Set pre-selections
func (m *Manager) Set(ctx context.Context, c ...Complete) error {
	return m.repository.Update(ctx, c...)
}*/
