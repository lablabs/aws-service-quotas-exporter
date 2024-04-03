package quota

import "github.com/aws/aws-sdk-go-v2/service/servicequotas"

type config struct {
	def    bool
	region string
}

func (c *config) options() func(*servicequotas.Options) {
	return func(ops *servicequotas.Options) {
		if c.region != "" {
			ops.Region = c.region
		}
	}
}

type Option func(c *config)

func WithDefault(def bool) Option {
	return func(c *config) {
		c.def = def
	}
}

func WithRegion(region string) Option {
	return func(c *config) {
		if region != "" {
			c.region = region
		}
	}
}
