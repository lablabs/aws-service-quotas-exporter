package quotas

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/servicequotas/types"
	"github.com/lablabs/aws-service-quotas-exporter/pkg/quota"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
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
	gvq := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: ns,
		Name:      "quota",
		Help:      "AWS service quota",
	}, []string{name, code, serviceCode})
	cl := Collector{
		log: log,
		qcl: qcl,
		gvq: gvq,
		cfg: cfg,
	}
	return &cl, nil
}

type Collector struct {
	log *logrus.Logger
	qcl Quota
	gvq *prometheus.GaugeVec
	cfg []Config
}

func (c *Collector) Register(r *prometheus.Registry) error {
	r.MustRegister(c.gvq)
	return nil
}

func (c *Collector) Collect(ctx context.Context, g *errgroup.Group) {
	for _, q := range c.cfg {
		g.Go(c.run(ctx, q))
	}
}

func (c *Collector) run(ctx context.Context, q Config) func() error {
	return func() error {
		res, err := c.qcl.GetQuota(ctx, q.ServiceCode, q.QuotaCode, quota.WithRegion(q.Region))
		if err != nil {
			return err
		}
		setMetric(c.gvq, res)
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
