package quota

import "github.com/aws/aws-sdk-go-v2/service/servicequotas"

type Options struct {
	*servicequotas.Options
}

type Option func(c *Options)

func WithRegion(region string) Option {
	return func(c *Options) {
		if region != "" {
			c.Region = region
		}
	}
}

func buildOptions(option ...Option) func(*servicequotas.Options) {
	return func(awsSq *servicequotas.Options) {
		op := Options{
			awsSq,
		}
		for _, o := range option {
			o(&op)
		}
	}
}
