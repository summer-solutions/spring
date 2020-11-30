package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetString(t *testing.T) {
	err := os.Setenv("SPRING_CONFIG_FILE", "../../config/web-api/config.test.yaml")
	assert.NoError(t, err)
	config, err := NewViperConfig("")

	assert.NoError(t, err)
	assert.Equal(t, "test_value", config.GetString("test_key"))
	assert.Equal(t, true, config.GetBool("test_bool"))
}

func TestViperConfig_GetMainPath(t *testing.T) {
	err := os.Setenv("SPRING_CONFIG_FILE", "../../config/web-api/config.test.yaml")
	assert.NoError(t, err)

	config, err := NewViperConfig("")
	assert.NoError(t, err)
	abs, err := filepath.Abs("../../config/")
	assert.NoError(t, err)
	assert.Equal(t, abs, config.GetMainPath())
}
