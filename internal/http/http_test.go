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

func TestNewHttp(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	address, close := StartHttp(t, ctx)
	defer func() {
		cancel()
		close()
	}()
	resp, err := http.Get("http://" + address + "/metrics")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func StartHttp(t *testing.T, ctx context.Context) (string, func()) {
	address := "0.0.0.0:8080"
	http, err := httpApi.NewHttp(test.DefaultLogger(), address, prometheus.NewRegistry())
	assert.NoError(t, err)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := http.Run(ctx)
		assert.NoError(t, err)
	}()
	return address, func() {
		wg.Wait()
	}
}
