package filters

import "github.com/schmichael/netdbg"

var (
	filterRegistry = map[string]FilterFactory{}
)

type FilterFactory func() netdbg.Filter

func RegisterFilter(name string, filter FilterFactory) {
	filterRegistry[name] = filter
}

func GetFilter(name string) FilterFactory {
	return filterRegistry[name]
}
