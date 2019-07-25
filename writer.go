package logger

import (
	"os"
)

type (
	filterWriter struct {
		file    *os.File
		filters []*filter
	}
)

func newFilterWriter(file *os.File, filters ...*filter) *filterWriter {
	if len(filters) == 0 {
		filters = append(filters, defaultFilters...)
	}

	return &filterWriter{
		file:    file,
		filters: filters,
	}
}

func (w *filterWriter) Write(p []byte) (n int, err error) {
	data := w.filter(p)

	return w.file.Write(data)
}

func (w *filterWriter) Sync() error {
	return w.file.Sync()
}

func (w *filterWriter) filter(data []byte) []byte {
	for _, filter := range w.filters {
		data = filter.reg.ReplaceAll(data, filter.repl)
	}

	return data
}
