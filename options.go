package outfmt

// Option configures rendering behavior.
type Option func(*renderOptions)

type renderOptions struct {
	jsonIndent string
	emptyValue string
}

func defaultRenderOptions() renderOptions {
	return renderOptions{
		jsonIndent: "  ",
		emptyValue: "-",
	}
}

// WithJSONIndent sets the indentation used by JSON output.
func WithJSONIndent(indent string) Option {
	return func(opts *renderOptions) {
		if indent != "" {
			opts.jsonIndent = indent
		}
	}
}

// WithEmptyValue sets the placeholder used for empty table cells.
func WithEmptyValue(value string) Option {
	return func(opts *renderOptions) {
		opts.emptyValue = value
	}
}
