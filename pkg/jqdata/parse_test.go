package jqdata

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func TestJsonData_Query(t *testing.T) {
	data := map[string]any{
		"v": 1,
		"s": "s",
		"l": []string{"1", "2", "3"},
	}
	j, err := json.Marshal(data)
	assert.NoError(t, err)
	type args struct {
		q string
	}
	tests := []struct {
		name    string
		j       func() JSONData
		args    args
		want    any
		wantErr bool
	}{
		{
			name: "Query OK",
			j: func() JSONData {
				jd, err := ParseRawJSON(j)
				assert.NoError(t, err)
				return jd
			},
			args: args{
				q: ".l | length",
			},
			want:    3,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
			defer cancel()
			jd := tt.j()
			got, err := jd.Query(ctx, tt.args.q)
			if (err != nil) != tt.wantErr {
				t.Errorf("Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Query() got = %v, want %v", got, tt.want)
			}
		})
	}
}
