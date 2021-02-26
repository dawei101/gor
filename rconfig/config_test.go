package rconfig

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegConfig(t *testing.T) {
	assert.Panics(t, func() { Reg("json", "test.json") })

	assert.Panics(t, func() { Reg("non", "non.yml") })

	assert.Panics(t, func() { Reg("wrongYml", "tests/wrong.yml") })

	assert.NotPanics(t, func() { Reg("defaultYml", "tests/config.yml") })
}

func TestGetConfig(t *testing.T) {
	Reg("default", "tests/config.yml")

	s := Get("default")

	type Mysql struct {
		DB   string
		Host string
		Port int
	}
	var mysql Mysql
	s.ValTo("mysql", &mysql)
	assert.Equal(t, Mysql{
		DB:   "default",
		Host: "127.0.0.1",
		Port: 3306,
	}, mysql)

	type Calendar struct {
		Days   []int
		Months map[int]string
	}
	var calendar Calendar
	s.ValTo("level1.calendar", &calendar)
	assert.Equal(t, Calendar{
		Days: []int{1, 2, 3, 4, 5},
		Months: map[int]string{
			1: "January",
			2: "February",
			3: "March",
		},
	}, calendar)
}
