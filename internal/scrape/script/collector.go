package script

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func NewCollector(log *logrus.Logger, cfg []Config, ns string) (*Collector, error) {

	tks := make([]task, 0)
	for _, c := range cfg {
		m := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: ns,
			Name:      c.Name,
			Help:      c.Help,
		}, c.LabelNames())
		tks = append(tks, task{
			m:   m,
			cfg: c,
		})
	}

	cl := Collector{
		log:   log,
		tasks: tks,
	}
	return &cl, nil
}

type Collector struct {
	log   *logrus.Logger
	tasks []task
}

func (c *Collector) Register(r *prometheus.Registry) error {
	for _, t := range c.tasks {
		r.MustRegister(t.m)
	}
	return nil
}

func (c *Collector) Collect(ctx context.Context, g *errgroup.Group) {
	for _, t := range c.tasks {
		g.Go(t.run(ctx))
	}
}

type task struct {
	m   *prometheus.GaugeVec
	cfg Config
}

func (t task) run(ctx context.Context) func() error {
	return func() error {
		data, err := Run(ctx, t.cfg)
		if err != nil {
			return err
		}
		for _, d := range data {
			t.m.With(d.Labels).Set(d.Value)
		}
		return nil
	}
}
