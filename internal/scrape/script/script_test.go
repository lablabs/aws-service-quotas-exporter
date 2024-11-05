package script_test

import (
	"context"
	"github.com/lablabs/aws-service-quotas-exporter/internal/scrape/script"
	"reflect"
	"testing"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name    string
		args    script.Config
		want    []script.Data
		wantErr bool
	}{
		{
			name: "Parsing command OK",
			args: script.Config{
				Name:   "metric_1",
				Help:   "",
				Script: "echo \"region=eu-central-1,cluster=eks-dev-1,type=dev,2\"",
			},
			want: []script.Data{
				{
					Value: 2,
					Labels: map[string]string{
						"region":  "eu-central-1",
						"cluster": "eks-dev-1",
						"type":    "dev",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid command",
			args: script.Config{
				Name:   "metric_1",
				Help:   "",
				Script: "invalid command",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancelCtx := context.WithCancel(context.Background())
			defer cancelCtx()
			got, err := script.Run(ctx, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() got = %v, error = %v, wantErr %v", got, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Run() got = %v, want %v", got, tt.want)
			}
		})
	}
}
