package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	t.Parallel()

	got, err := GetConfig()
	assert.NoError(t, err)

	assert.Equal(t, got, &Config{
		ListenPort:  8082,
		MetricsPort: 9090,
	})
}
