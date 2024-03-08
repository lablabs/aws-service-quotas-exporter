package exporter_test

import (
	"context"
	"errors"
	"github.com/lablabs/aws-service-quotas-exporter/internal/exporter"
	"github.com/lablabs/aws-service-quotas-exporter/test"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
	"testing"
	"time"
)

func TestExporter_Run(t *testing.T) {
	type fields struct {
		log *logrus.Logger
		cls []exporter.Collector
		ops []exporter.Option
		ctx func() (context.Context, func())
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Exporter OK",
			fields: fields{
				log: test.DefaultLogger(),
				cls: []exporter.Collector{&testCollector{}},
				ops: []exporter.Option{},
				ctx: func() (context.Context, func()) {
					return context.WithTimeout(context.Background(), time.Second*1)
				},
			},
			wantErr: false,
		},
		{
			name: "Exporter timeout",
			fields: fields{
				log: test.DefaultLogger(),
				cls: []exporter.Collector{&testCollector{err: errors.New("timeout")}},
				ops: []exporter.Option{},
				ctx: func() (context.Context, func()) {
					return context.WithTimeout(context.Background(), time.Second*1)
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := exporter.NewExporter(tt.fields.log, tt.fields.cls)
			assert.NoError(t, err)
			ctx, cancel := tt.fields.ctx()
			defer cancel()
			if err := e.Run(ctx); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}

type testCollector struct {
	err error
}

func (t *testCollector) Register(_ *prometheus.Registry) error {
	return nil
}

func (t *testCollector) Collect(_ context.Context, g *errgroup.Group) {
	g.Go(func() error {
		return t.err
	})
}
