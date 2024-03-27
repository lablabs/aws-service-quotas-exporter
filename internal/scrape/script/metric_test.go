package script_test

import (
	"github.com/lablabs/aws-service-quotas-exporter/internal/scrape/script"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestParser_ParseMetric(t *testing.T) {
	tests := []struct {
		name    string
		metric  string
		want    script.Data
		wantErr bool
	}{
		{
			name:   "Parse Metric OK",
			metric: "a=1,b=2,c='3',d='t e a',y=\"aa bbb ccc\",4",
			want: script.Data{
				Value: 4,
				Labels: map[string]string{
					"a": "1",
					"b": "2",
					"c": "3",
					"d": "t e a",
					"y": "aa bbb ccc",
				},
			},
			wantErr: false,
		},
		{
			name:   "Value not valid float64",
			metric: "a=1,b=2,c='3',d='t e a',y=\"aa bbb ccc\",not_valid",
			want: script.Data{
				Value: 4,
				Labels: map[string]string{
					"a": "1",
					"b": "2",
					"c": "3",
					"d": "t e a",
					"y": "aa bbb ccc",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := script.NewParser()
			got, err := p.ParseMetric(tt.metric)
			if tt.wantErr {
				assert.Errorf(t, err, "ParseMetric() expect error")
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseMetric() got = %v, want %v", got, tt.want)
			}
		})
	}
}
