package exporter

import "time"

type config struct {
	interval time.Duration
	timeout  time.Duration
}

type Option func(c *config) error

func WithInterval(i time.Duration) Option {
	return func(c *config) error {
		if i.Nanoseconds() != 0 {
			c.interval = i
		}
		return nil
	}
}

func WithTimeout(t time.Duration) Option {
	return func(c *config) error {
		if t.Nanoseconds() != 0 {
			c.timeout = t
		}
		return nil
	}
}
