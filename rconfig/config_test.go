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

	var m mysql
	s.ValTo("mysql", &m)
	assert.Equal(t, expectedMysql, m)

	var c calendar
	s.ValTo("level1.calendar", &c)
	assert.Equal(t, expectedCalendar, c)

	var h human
	s.ValTo("human", &h)
	assert.Equal(t, expectedHuman, h)
}

func TestDefConf(t *testing.T) {
	Reg("default", "tests/config.yml")

	var m mysql
	DefConf().ValTo("mysql", &m)
	assert.Equal(t, expectedMysql, m)

	var c calendar
	DefConf().ValTo("level1.calendar", &c)
	assert.Equal(t, expectedCalendar, c)

	var h human
	DefConf().ValTo("human", &h)
	assert.Equal(t, expectedHuman, h)
}

type mysql struct {
	DB   string
	Host string
	Port int
}

type calendar struct {
	Days   []int
	Months map[int]string
}

type fruit struct {
	ID   int
	Name string
}
type human struct {
	Name  string
	Likes []fruit
}

var expectedMysql = mysql{
	DB:   "default",
	Host: "127.0.0.1",
	Port: 3306,
}

var expectedCalendar = calendar{
	Days: []int{1, 2, 3, 4, 5},
	Months: map[int]string{
		1: "January",
		2: "February",
		3: "March",
	},
}

var expectedHuman = human{
	Name: "mary",
	Likes: []fruit{
		{
			ID:   1,
			Name: "mango",
		},
		{
			ID:   2,
			Name: "apple",
		},
		{
			ID:   3,
			Name: "banana",
		},
	},
}
