package filters

var (
	filterRegistry = map[string]FilterFactory{}
)

// FilterFactory is the New function Filters must implement to be used.
type FilterFactory func(writerIn chan []byte, readerIn chan []byte) (f Filter, writerOut chan []byte, readerOut chan []byte)

// RegisterFilter registers a Filter by a name.
func RegisterFilter(name string, filter FilterFactory) {
	filterRegistry[name] = filter
}

// GetFilter returns a FilterFactory for a given name or nil if no Filter has
// been registered for that name.
func GetFilter(name string) FilterFactory {
	return filterRegistry[name]
}
