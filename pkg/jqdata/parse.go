package jqdata

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/itchyny/gojq"
)

type JSONData struct {
	d map[string]any
}

func ParseRawJSON(data []byte) (JSONData, error) {
	d := make(map[string]any)
	err := json.Unmarshal(data, &d)
	if err != nil {
		return JSONData{}, fmt.Errorf("unable parse JSON: %w", err)
	}
	return JSONData{d: d}, nil
}

func (j JSONData) Query(ctx context.Context, q string) (any, error) {
	qr, err := gojq.Parse(q)
	if err != nil {
		return nil, err
	}
	it := qr.RunWithContext(ctx, j.d)
	for {
		v, ok := it.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			return nil, err
		}
		return v, nil
	}
	return nil, fmt.Errorf("empty data for query: %v", q)
}
