package logger

import (
	"regexp"
	"strings"
)

type (
	filter struct {
		reg  *regexp.Regexp
		repl []byte
	}
)

var (
	defaultFilters []*filter
)

func init() {
	defaultFilters = append(
		defaultFilters,
		newFilter(`"password":(\s*)".*?"`, `"password":$1"*"`),
		newFilter(`password:(\s*).*?\S*`, `password:$1*`),
		newFilter(`\\"password\\":(\s*)\\".*?\\"`, `\"password\":$1\"*\"`),
		newFilter(`password=\w*&`, `password=*&`),
		newFilter(`password=\w*\S`, `password=*`),
	)
}

func newFilterBySlice(specs []string) []*filter {
	filters := []*filter{}

	for _, spec := range specs {
		words := strings.Split(spec, "==>")
		if len(words) != 2 {
			continue
		} else {
			Warnf("newFilter: invalid spec %s", spec)
		}

		filters = append(filters, newFilter(words[0], words[1]))
	}

	return filters
}

func newFilter(rule, repl string) *filter {
	return &filter{
		reg:  regexp.MustCompile(rule),
		repl: []byte(repl),
	}
}

// filterDefault
func filterDefault(data []byte) []byte {
	for _, filter := range defaultFilters {
		data = filter.reg.ReplaceAll(data, filter.repl)
	}

	return data
}
