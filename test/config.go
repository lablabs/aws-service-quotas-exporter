package test

import (
	_ "embed"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

//go:embed configs/zone.yaml
var configEip string

func TmpConfigMetricFile(t *testing.T) string {
	f, err := os.CreateTemp(os.TempDir(), "exporter")
	assert.NoError(t, err)
	_, err = f.WriteString(configEip)
	assert.NoError(t, err)
	return f.Name()
}
