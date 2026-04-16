package outfmt

import "fmt"

// Format identifies the output renderer used by Render and Marshal.
type Format string

const (
	// Table renders data as an aligned text table for terminal reading.
	Table Format = "table"
	// YAML renders data using YAML serialization while preserving nesting.
	YAML Format = "yaml"
	// JSON renders data as pretty-printed JSON.
	JSON Format = "json"
)

// Validate returns an error when the format is not supported.
func (f Format) Validate() error {
	switch f {
	case Table, YAML, JSON:
		return nil
	default:
		return fmt.Errorf("outfmt: unsupported format %q", f)
	}
}
