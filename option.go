package pi

type Options struct {
	group string
	chain []Middleware
}

type Option func(*Options)

func Chain(chain ...Middleware) Option {
	return func(o *Options) {
		o.chain = chain
	}
}

func Group(group string) Option {
	return func(o *Options) {
		o.group = group
	}
}
