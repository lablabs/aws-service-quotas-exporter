package http_test

import (
	"context"
	httpApi "github.com/lablabs/aws-service-quotas-exporter/internal/http"
	"github.com/lablabs/aws-service-quotas-exporter/test"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"net/http"
	"sync"
	"testing"
)

const (
	address = "0.0.0.0:8080"
)

func TestNewHttp(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	closeHTTP := StartHTTP(ctx, t)
	defer func() {
		cancel()
		closeHTTP()
	}()
	resp, err := http.Get("http://" + address + "/metrics")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func StartHTTP(ctx context.Context, t *testing.T) func() {
	http, err := httpApi.NewHTTP(test.DefaultLogger(), address, prometheus.NewRegistry())
	assert.NoError(t, err)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := http.Run(ctx)
		assert.NoError(t, err)
	}()
	return func() {
		wg.Wait()
	}
}
