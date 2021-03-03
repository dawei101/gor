package rlog

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/dawei101/gor/rconfig"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	rconfig.Reg("default", "tests/config.yml")
	assert.Nil(t, loadConfig())
}

func TestRotateAllLogs(t *testing.T) {
	rconfig.Reg("default", "tests/config.yml")

	loadConfig()
	log := GetLog(defName)
	_, logfn := filepath.Split(log.filename)
	prefix := strings.TrimSuffix(logfn, ".log") + "-" + time.Now().Format("2006-01-02")

	fmt.Println("logfile name prefix:", prefix)

	filepath.Walk(config.Path, func(path string, info os.FileInfo, err error) error {
		if strings.TrimPrefix(info.Name(), prefix) != info.Name() {
			os.Remove(path)
		}
		return nil
	})

	Info(context.Background(), "-----")
	rotateAllLogs()

	ok := false
	filepath.Walk(config.Path, func(path string, info os.FileInfo, err error) error {
		if strings.TrimPrefix(info.Name(), prefix) != info.Name() {
			ok = true
		}
		return nil
	})

	assert.Equal(t, true, ok, "should find rotated log file")
}

func TestLogFunc(t *testing.T) {
	rconfig.Reg("default", "tests/config.yml")

	loadConfig()
	log := GetLog(defName)

	logLevel = LEVEL_DEBUG

	log.Info(context.Background(), ">>>")
	log.Warning(context.Background(), ">>>")
	log.Debug(context.Background(), ">>>")
	log.Error(context.Background(), ">>>")

	dat, _ := ioutil.ReadFile(log.filename)
	s := string(dat)
	assert.Equal(t, 1, strings.Count(s, "Info"), "one info record")
	assert.Equal(t, 1, strings.Count(s, "Warning"), "one warning record")
	assert.Equal(t, 1, strings.Count(s, "Error"), "one error record")
	assert.Equal(t, 1, strings.Count(s, "Debug"), "one debug record")

	os.Remove(log.filename)
}

func TestLogSerialNum(t *testing.T) {

	myc := counter{}
	myc.Rise()
	myc.Rise()
	myc.Rise()
	c := myc.Rise()
	assert.Equal(t, int32(4), c, "counter should be 4")

	r, _ := http.NewRequest(http.MethodPost, "/test", nil)
	r = prepareRequest(r)

	idx0 := logSerialNum(r.Context())

	idx1 := logSerialNum(r.Context())

	assert.Equal(t, idx0+1, idx1, "counter rise 1 once")
}
