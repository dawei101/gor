package rlog

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/dawei101/gor/rconfig"
)

var config struct {
	Path       string `yaml:"path"`
	Level      string `yaml:"level"`
	MaxMB      int    `yaml:"maxDB"`
	MaxDays    int    `yaml:"MaxDays"`
	MaxBackups int    `yaml:"MaxBackups"`
}

func loadConfig() error {

	rconfig.DefConf().ValTo("rlog", &config)

	if config.Path == "" {
		return errors.New("Log config is not corrent")
	}
	if _, err := os.Stat(config.Path); err != nil {
		return errors.New("Log config(log.path) is not correct:" + config.Path)
	}

	logLevel = logLevelFromString(config.Level)
	deflog := New(defName)
	logs.Store(defName, deflog)

	Debug = deflog.Debug
	Warning = deflog.Warning
	Error = deflog.Error
	Info = deflog.Info

	return nil
}

func logLevelFromString(level string) LogLevel {
	switch strings.ToLower(level) {
	case "debug":
		return LEVEL_DEBUG
	case "info":
		return LEVEL_INFO
	case "warning", "warn":
		return LEVEL_WARNING
	case "error":
		return LEVEL_ERROR
	default:
		return LEVEL_INFO
	}
}

func midnightRotate() {
	oneDay := 24 * time.Hour
	t := time.Now()
	midn := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
	d := midn.Sub(t)
	if d < 0 {
		midn = midn.Add(oneDay)
		d = midn.Sub(t)
	}
	for {
		time.Sleep(d)
		d = oneDay
		go rotateAllLogs()
	}
}

func rotateAllLogs() {
	logs.Range(func(key, value interface{}) bool {
		log := value.(*Log)
		return log.writer.Rotate() != nil
	})
}
