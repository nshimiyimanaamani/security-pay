package stats

// Chart is a stats aggregate
type Chart struct {
	Label string            `json:"label,omitempty"`
	Data  map[string]uint64 `json:"data,omitempty"`
}
