package outfmt

import "fmt"

type Format string

const (
	Table Format = "table"
	YAML  Format = "yaml"
	JSON  Format = "json"
)

func (f Format) Validate() error {
	switch f {
	case Table, YAML, JSON:
		return nil
	default:
		return fmt.Errorf("outfmt: unsupported format %q", f)
	}
}
