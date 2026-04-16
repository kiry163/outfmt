package outfmt

import (
	"fmt"
	"io"
	"strings"
)

func writeTable(w io.Writer, data tableRows, opts renderOptions) error {
	if len(data.headers) == 0 {
		_, err := io.WriteString(w, "(no data)\n")
		return err
	}

	widths := make([]int, len(data.headers))
	for i, header := range data.headers {
		widths[i] = len(header)
	}

	normalizedRows := make([][]string, 0, len(data.rows))
	for _, row := range data.rows {
		normalized := make([]string, len(data.headers))
		for i := range data.headers {
			cell := ""
			if i < len(row) {
				cell = row[i]
			}
			if cell == "" {
				cell = opts.emptyValue
			}
			normalized[i] = cell
			if len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
		normalizedRows = append(normalizedRows, normalized)
	}

	if _, err := io.WriteString(w, formatRow(data.headers, widths)); err != nil {
		return err
	}
	if _, err := io.WriteString(w, formatSeparator(widths)); err != nil {
		return err
	}
	for _, row := range normalizedRows {
		if _, err := io.WriteString(w, formatRow(row, widths)); err != nil {
			return err
		}
	}
	return nil
}

func formatRow(values []string, widths []int) string {
	parts := make([]string, 0, len(values))
	for i, value := range values {
		parts = append(parts, padRight(value, widths[i]))
	}
	return strings.Join(parts, "  ") + "\n"
}

func formatSeparator(widths []int) string {
	parts := make([]string, 0, len(widths))
	for _, width := range widths {
		parts = append(parts, strings.Repeat("-", width))
	}
	return strings.Join(parts, "  ") + "\n"
}

func padRight(value string, width int) string {
	if len(value) >= width {
		return value
	}
	return value + strings.Repeat(" ", width-len(value))
}

func mustWriteTable(data any, opts ...Option) string {
	raw, err := Marshal(data, Table, opts...)
	if err != nil {
		panic(fmt.Sprintf("unexpected table error: %v", err))
	}
	return string(raw)
}
