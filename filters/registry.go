package filters

import "github.com/schmichael/netdbg"

var (
	filterRegistry = map[string]FilterFactory{}
)

// FilterFactory is the New function Filters must implement to be used.
type FilterFactory func() netdbg.Filter

// RegisterFilter registers a Filter by a name.
func RegisterFilter(name string, filter FilterFactory) {
	filterRegistry[name] = filter
}

// GetFilter returns a FilterFactory for a given name or nil if no Filter has
// been registered for that name.
func GetFilter(name string) FilterFactory {
	return filterRegistry[name]
}
