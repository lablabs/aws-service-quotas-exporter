package script

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"sync"
)

func NewCollector(log *logrus.Logger, cfg []Config, ns string) (*Collector, error) {
	cl := Collector{
		log:   log,
		ns:    ns,
		cfg:   cfg,
		tasks: make([]task, 0),
	}
	return &cl, nil
}

type Collector struct {
	log       *logrus.Logger
	once      sync.Once
	err       error
	ns        string
	cfg       []Config
	tasks     []task
	tasksLock sync.Mutex
}

func (c *Collector) Register(ctx context.Context, r *prometheus.Registry) error {
	c.once.Do(func() {
		c.log.Debugf("start registering script metrics")
		g, ctx := errgroup.WithContext(ctx)
		for _, cf := range c.cfg {
			config := cf
			g.Go(func() error {
				data, err := Run(ctx, config)
				if err != nil {
					c.log.Errorf("unable to run command: %s, %v", config.Script, err)
					return err
				}
				if len(data) > 0 {
					lbs := data[0].LabelNames()
					m := prometheus.NewGaugeVec(prometheus.GaugeOpts{
						Namespace: c.ns,
						Name:      config.Name,
						Help:      config.Help,
					}, lbs)
					r.MustRegister(m)
					t := task{
						m:   m,
						cfg: config,
					}
					c.addTask(t)
					for _, d := range data {
						t.m.With(d.Labels).Set(d.Value)
					}
				}
				return nil
			})
		}
		c.err = g.Wait()
	})
	return c.err
}

func (c *Collector) Collect(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)
	for _, t := range c.tasks {
		g.Go(t.run(ctx))
	}
	err := g.Wait()
	return err
}

func (c *Collector) addTask(t task) {
	c.tasksLock.Lock()
	defer c.tasksLock.Unlock()
	c.tasks = append(c.tasks, t)
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
