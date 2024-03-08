package script

import (
	"context"
	"reflect"
	"testing"
	"time"
)

func TestScrapper_Run(t *testing.T) {
	type fields struct {
		cfg Config
	}
	tests := []struct {
		name    string
		fields  fields
		want    []Data
		wantErr bool
	}{
		{
			name: "Scrape items OK",
			fields: fields{
				cfg: Config{
					Command: "echo {\"items\":[{\"v\": 1,\"name\":\"a\"},{\"v\":2,\"name\":\"b\"}]}",
					List:    ".items",
					Value:   ".v",
					Labels: []Label{
						{
							Name:    "name",
							JqValue: ".name",
						},
					},
				},
			},
			want: []Data{
				{
					Value: 1,
					Labels: map[string]string{
						"name": "a",
					},
				},
				{
					Value: 2,
					Labels: map[string]string{
						"name": "b",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Scrape one value OK",
			fields: fields{
				cfg: Config{
					Command: "echo {\"v\": 1,\"name\":\"a\"}",
					Value:   ".v",
					Labels: []Label{
						{
							Name:    "name",
							JqValue: ".name",
						},
					},
				},
			},
			want: []Data{
				{
					Value: 1,
					Labels: map[string]string{
						"name": "a",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Scrape cmd Json Error",
			fields: fields{
				cfg: Config{
					Command: "echo not json",
					Value:   ".v",
					Labels: []Label{
						{
							Name:    "name",
							JqValue: ".name",
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			got, err := Run(ctx, tt.fields.cfg)
			cancel()
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Run() got = %v, want %v", got, tt.want)
			}
		})
	}
}
