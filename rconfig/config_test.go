package rconfig

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegConfig(t *testing.T) {
	assert.Panics(t, func() { RegConfig("json", "test.json") })

	assert.Panics(t, func() { RegConfig("non", "non.yml") })

	assert.Panics(t, func() { RegConfig("wrongYml", "tests/wrong.yml") })

	assert.NotPanics(t, func() { RegConfig("defaultYml", "tests/config.yml") })
}

func TestGetConfig(t *testing.T) {
	RegConfig("default", "tests/config.yml")

	s := GetConfig("default")

	type Mysql struct {
		DB   string
		Host string
		Port int
	}
	var mysql Mysql
	s.ValueAssignTo("mysql", &mysql, "")
	assert.Equal(t, mysql, Mysql{
		DB:   "default",
		Host: "127.0.0.1",
		Port: 3306,
	})

	type Calendar struct {
		Days   []int
		Months map[int]string
	}
	var calendar Calendar
	s.ValueAssignTo("calendar", &calendar, "")
	assert.Equal(t, calendar, Calendar{
		Days: []int{1, 2, 3, 4, 5},
		Months: map[int]string{
			1: "January",
			2: "February",
			3: "March",
		},
	})
}
