package htmlx

type Config struct {
	// contains filtered or unexported fields

	async bool
}

type Option func(c *Config)

// Async specify whether the parsing is done asynchronously or not. it is synchronous by default
func Async(allow bool) Option {
	return func(c *Config) {
		c.async = allow
	}
}
