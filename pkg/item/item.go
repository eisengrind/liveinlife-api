package item

// ConfigEntry is a config entry of a virtual GTA5 item
type ConfigEntry struct {
	Weight    float64 `json:"weight"`
	MaxSubset float64 `json:"max_subset"`
}

// Config of available virtual GTA5 items
type Config map[string]*ConfigEntry
