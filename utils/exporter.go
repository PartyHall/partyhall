package utils

type Exporter interface {
	Export() (map[string]any, error)
}
