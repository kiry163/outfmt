package outfmt

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

func WithJSONIndent(indent string) Option {
	return func(opts *renderOptions) {
		if indent != "" {
			opts.jsonIndent = indent
		}
	}
}

func WithEmptyValue(value string) Option {
	return func(opts *renderOptions) {
		opts.emptyValue = value
	}
}
