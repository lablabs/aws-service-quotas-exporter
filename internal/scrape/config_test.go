package scrape

import (
	"testing"
	"time"
)

func TestScrape_Validate(t *testing.T) {
	type fields struct {
		Interval time.Duration
		Timeout  time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "Scrape config not set",
			fields:  fields{},
			wantErr: false,
		},
		{
			name: "Scrape config OK",
			fields: fields{
				Interval: time.Minute + time.Second*5,
				Timeout:  time.Second * 6,
			},
			wantErr: false,
		},
		{
			name: "Scrape config error",
			fields: fields{
				Interval: time.Minute - time.Second,
				Timeout:  time.Second * 4,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scrape{
				Interval: tt.fields.Interval,
				Timeout:  tt.fields.Timeout,
			}
			if err := s.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
