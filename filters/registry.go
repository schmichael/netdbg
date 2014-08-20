package filters

var (
	filterRegistry = map[string]FilterFactory{}
)

// RegisterFilter registers a Filter by a name.
func RegisterFilter(name string, filter FilterFactory) {
	filterRegistry[name] = filter
}

// GetFilter returns a FilterFactory for a given name or nil if no Filter has
// been registered for that name.
func GetFilter(name string) FilterFactory {
	return filterRegistry[name]
}
