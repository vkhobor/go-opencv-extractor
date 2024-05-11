package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMustNewDefaultViper_NoPanic(t *testing.T) {
	assert.NotPanics(t, func() {
		MustNewDefaultViperConfig()
	})
}
