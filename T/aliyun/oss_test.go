package aliyun

import (
	"fmt"
	"path"
	"runtime"
	"testing"

	"github.com/dawei101/gor"
)

func TestOss(t *testing.T) {

	_, filename, _, _ := runtime.Caller(0)
	cf := path.Join(path.Dir(filename), "../test_config.yml")
	fmt.Printf("config file set to: %s\n", cf)
	rlib.SetConfigFile(cf)
	rlib.DefaultRooboConfig()

	b, err := DefaultOssBucket()
	if err != nil {
		t.Errorf(err.Error())
	}

	fmt.Println("put file:" + cf)
	err = b.PutObjectFromFile("test_config.yml", cf)
	if err != nil {
		t.Errorf(err.Error())
	}
}
