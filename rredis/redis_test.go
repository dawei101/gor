package rredis

import (
	"testing"

	"github.com/dawei101/gor/rconfig"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	rconfig.Reg(DefName, "tests/config.yml")
	assert.Panics(t, func() { loadConfig() })

	assert.Panics(t, func() { Get("noneredis") })
}
