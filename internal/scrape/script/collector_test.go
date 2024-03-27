package script_test

import (
	"context"
	"github.com/lablabs/aws-service-quotas-exporter/internal/scrape/script"
	"github.com/lablabs/aws-service-quotas-exporter/test"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewCollector(t *testing.T) {

	cfg := []script.Config{
		{
			Name:   "metric_a",
			Help:   "metric b",
			Script: "echo \"name=n,cluster=1,1\"",
		},
		{
			Name:   "metric_b",
			Help:   "metric b",
			Script: "echo \"name=n,cluster=2,1\"",
		},
	}

	cl, err := script.NewCollector(test.DefaultLogger(), cfg, "ns")
	assert.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	r := prometheus.NewRegistry()
	err = cl.Register(ctx, r)
	assert.NoError(t, err)
	err = cl.Register(ctx, r)
	assert.NoError(t, err)
	err = cl.Collect(ctx)
	assert.NoError(t, err)

}
