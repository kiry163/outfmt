package outfmt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

func Render(w io.Writer, data any, format Format, opts ...Option) error {
	if w == nil {
		return fmt.Errorf("outfmt: writer is nil")
	}

	raw, err := Marshal(data, format, opts...)
	if err != nil {
		return err
	}

	_, err = w.Write(raw)
	return err
}

func Marshal(data any, format Format, opts ...Option) ([]byte, error) {
	if err := format.Validate(); err != nil {
		return nil, err
	}

	config := defaultRenderOptions()
	for _, opt := range opts {
		if opt != nil {
			opt(&config)
		}
	}

	switch format {
	case JSON:
		return marshalJSON(data, config)
	case YAML:
		return marshalYAML(data)
	case Table:
		return marshalTable(data, config)
	default:
		return nil, fmt.Errorf("outfmt: unsupported format %q", format)
	}
}

func marshalJSON(data any, opts renderOptions) ([]byte, error) {
	raw, err := json.MarshalIndent(data, "", opts.jsonIndent)
	if err != nil {
		return nil, fmt.Errorf("outfmt: marshal json: %w", err)
	}
	return append(raw, '\n'), nil
}

func marshalYAML(data any) ([]byte, error) {
	raw, err := yaml.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("outfmt: marshal yaml: %w", err)
	}
	return raw, nil
}

func marshalTable(data any, opts renderOptions) ([]byte, error) {
	rows, err := extractRows(data)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := writeTable(&buf, rows, opts); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
