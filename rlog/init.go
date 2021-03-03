// +build !testing

package rlog

func init() {
	if err := loadConfig(); err != nil {
		panic(err)
	}
	go midnightRotate()
}
