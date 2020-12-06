package spring

import (
	"testing"

	"github.com/sarulabs/di"
	"github.com/tj/assert"
)

func TestFlags(t *testing.T) {
	r := New("test_script").RegisterDIService()
	testService := &ServiceDefinition{
		Name:   "test_service",
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			return "hello", nil
		},
		Flags: func(registry *FlagsRegistry) {
			registry.Bool("test-bool-false", false, "")
			registry.Bool("test-bool-true", true, "")
			registry.String("test-string", "default string", "")
		},
	}
	r.RegisterDIService(testService).Build()

	assert.Equal(t, false, DIC().App().Flags().Bool("test-bool-false"))
	assert.Equal(t, true, DIC().App().Flags().Bool("test-bool-true"))
	assert.Equal(t, "default string", DIC().App().Flags().String("test-string"))
	assert.Equal(t, false, DIC().App().Flags().Bool("missing-bool"))
	assert.Equal(t, "", DIC().App().Flags().String("missing-string"))
}
