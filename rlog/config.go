package rlog

import (
	"github.com/dawei101/gor/rconfig"
)

var LogConfig struct {
	Path       string `json:"path" yaml:"path"`
	Level      string `json:"level" yaml:"level"`
	MaxMB      int    `json:"maxMB" yaml:"maxDB"`
	MaxDays    int    `json:"maxDays" yaml:"MaxDays"`
	MaxBackups int    `json:"maxBackups" yaml:"MaxBackups"`
	RotateTime string `json:"rotateTime" yaml:"RotateTime"`
}

func init() {
	rconfig.ValueMustAssignTo("rconfig", &LogConfig)

	c := LogConfig

	if c.Log == nil || len(c.Log.Path) == 0 {
		panic("Log config is not corrent")
	}
	if _, err := os.Stat(c.Log.Path); err != nil {
		panic("Log config(log.path) is not correct:" + c.Log.Path)
	}

	execf := filepath.Base(os.Args[0])

	reqlogf := path.Join(c.Log.Path, execf+".request.log")
	apilogf := path.Join(c.Log.Path, execf+".api.log")
	applogf := path.Join(c.Log.Path, execf+".app.log")

	loglv := LogLevelFromString(c.Log.Level)
	reqLog = NewFileLoggerHolder(reqlogf, loglv, c.Log.MaxMB, c.Log.MaxDays, c.Log.MaxBackups)
	apiLog = NewFileLoggerHolder(apilogf, loglv, c.Log.MaxMB, c.Log.MaxDays, c.Log.MaxBackups)
	appLog = NewFileLoggerHolder(applogf, loglv, c.Log.MaxMB, c.Log.MaxDays, c.Log.MaxBackups)
	SetDefaultLog(appLog)

	rotateTime := c.Log.RotateTime
	hour := 23
	min := 59
	sec := 59
	if rotateTime != "" {
		rotateSlice := strings.Split(rotateTime, ":")
		if len(rotateSlice) != 3 {
			panic("config rotate time should be like 23:59:59")
		}
		var err1, err2, err3 error
		hour, err1 = strconv.Atoi(rotateSlice[0])
		min, err2 = strconv.Atoi(rotateSlice[1])
		sec, err3 = strconv.Atoi(rotateSlice[2])
		if err1 != nil || err2 != nil || err3 != nil {
			panic("config rotate time should be like 23:59:59")
		}
		if hour < 0 || hour > 23 || min < 0 || min > 59 || sec < 0 || sec > 59 {
			panic("config rotate time invalidate,it should be like 23:59:59")
		}
	}

	go run_rotate_log(hour, min, sec)

}
