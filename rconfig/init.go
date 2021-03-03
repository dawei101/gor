// +build !testing

package rconfig

import (
	"os"
)

func init() {
	_, err := os.Stat(DefFile)
	if err != nil {
		panic(err)
	}

	if err = Reg(DefName, DefFile); err != nil {
		panic(err)
	}
}
