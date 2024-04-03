package quotas

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/servicequotas/types"
	"github.com/lablabs/aws-service-quotas-exporter/pkg/quota"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"sync"
)

const (
	name        = "name"
	code        = "code"
	serviceCode = "service_code"
)

type Quota interface {
	GetQuota(ctx context.Context, serviceCode string, quotaCode string, options ...quota.Option) (*types.ServiceQuota, error)
}

func NewCollector(log *logrus.Logger, cfg []Config, ns string, qcl Quota) (*Collector, error) {
	cl := Collector{
		log:   log,
		qcl:   qcl,
		cfg:   cfg,
		ns:    ns,
		tasks: make([]task, 0),
	}
	return &cl, nil
}

type Collector struct {
	log       *logrus.Logger
	qcl       Quota
	once      sync.Once
	err       error
	ns        string
	cfg       []Config
	tasks     []task
	tasksLock sync.Mutex
}

func (c *Collector) Register(ctx context.Context, r *prometheus.Registry) error {
	c.once.Do(func() {
		c.log.Debugf("start registering quota metrics")
		gvq := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: c.ns,
			Name:      "quota",
			Help:      "AWS service quota",
		}, []string{name, code, serviceCode})
		r.MustRegister(gvq)
		g, ctx := errgroup.WithContext(ctx)
		for _, cf := range c.cfg {
			qc := cf
			g.Go(func() error {
				res, err := c.qcl.GetQuota(ctx, qc.ServiceCode, qc.QuotaCode, quota.WithDefault(qc.Default), quota.WithRegion(qc.Region))
				if err != nil {
					return err
				}
				t := task{
					m:   gvq,
					cfg: qc,
				}
				c.addTask(t)
				setMetric(t.m, res)
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
		ts := t
		g.Go(ts.run(ctx, c.qcl))
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

func (t task) run(ctx context.Context, c Quota) func() error {
	return func() error {
		res, err := c.GetQuota(ctx, t.cfg.ServiceCode, t.cfg.QuotaCode, quota.WithDefault(t.cfg.Default), quota.WithRegion(t.cfg.Region))
		if err != nil {
			return err
		}
		setMetric(t.m, res)
		return nil
	}
}

func setMetric(gc *prometheus.GaugeVec, q *types.ServiceQuota) {
	gc.With(prometheus.Labels{
		name:        aws.ToString(q.QuotaName),
		code:        aws.ToString(q.QuotaCode),
		serviceCode: aws.ToString(q.ServiceCode),
	}).Set(aws.ToFloat64(q.Value))
}
