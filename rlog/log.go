package rlog

/*

本文件内的封装，日志将打入三个文件：
1. 请求     PATH/request.log
2. app      PATH/app.log
3. 内部调用 PATH/api.log

在开应用时，打印日志只需要调用
rlib.Debug
rlib.Info
rlib.Warning
rlib.Error

*/

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"sync/atomic"
	"time"

	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

type ContextLogFunc func(ctx context.Context, v ...interface{})
type LogLevel int

const (
	RequestIdKey = "*req*"
)

const (
	LEVEL_DEBUG LogLevel = iota - 1
	LEVEL_INFO
	LEVEL_WARNING
	LEVEL_ERROR
)

//
// string型日志级别快速转换为 `LogLevel`
//
func LogLevelFromString(level string) LogLevel {
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

// 日志快捷方法: `Debug` `Info` `Warning` `Error`
//
// 要切换日志输出目录，使用 `SetDefaultLog()` 设置新的RLog实例
//
var (
	glocation *time.Location
	glog      *RLog
	codeinfo  = true
	Level     = LEVEL_INFO
)

func init() {
	glocation, _ = time.LoadLocation("Asia/Chongqing")
	glog = newStdLog()
	go run_rotate_log()
}

type RLog struct {
	debugLogger   *myLogger
	infoLogger    *myLogger
	warningLogger *myLogger
	errorLogger   *myLogger
	logfile       string
	Debug         ContextLogFunc
	Info          ContextLogFunc
	Warning       ContextLogFunc
	Error         ContextLogFunc
}

func newRLog(logfile string, dlogger, ilogger, wlogger, elogger *myLogger) *RLog {
	return &RLog{
		debugLogger:   dlogger,
		infoLogger:    ilogger,
		warningLogger: wlogger,
		errorLogger:   elogger,
		logfile:       logfile,
		Debug:         dlogger.printlnX,
		Info:          ilogger.printlnX,
		Warning:       wlogger.printlnX,
		Error:         elogger.printlnX,
	}
}

func newStdLog() *RLog {
	c := uint64(0)
	return newRLog("/dev/stdout",
		newLogger(os.Stdout, "Debug:", LEVEL_DEBUG, &c),
		newLogger(os.Stdout, "Info:", LEVEL_INFO, &c),
		newLogger(os.Stdout, "Warning:", LEVEL_WARNING, &c),
		newLogger(os.Stderr, "Error:", LEVEL_ERROR, &c))
}

func NewFileLog(file string, loglv LogLevel, maxMB int, maxdays int, maxbackups int) *RLog {
	Level = loglv
	logfile := file
	_, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
		os.Exit(1)
	}
	writer := &lumberjack.Logger{
		Filename:   logfile,
		MaxSize:    maxMB,      // megabytes after which new file is created
		MaxBackups: maxbackups, // number of backups
		MaxAge:     maxdays,    //days
	}
	go care_rotate_log(writer)

	c := uint64(0)
	return newRLog(logfile,
		newLogger(writer, "Debug:", LEVEL_DEBUG, &c),
		newLogger(writer, "Info:", LEVEL_INFO, &c),
		newLogger(writer, "Warning:", LEVEL_WARNING, &c),
		newLogger(writer, "Error:", LEVEL_ERROR, &c))
}

//
// 关闭显示代码与行信息
//
func LogDisableCodeinfo() {
	codeinfo = false
}

//
// 设置日志打印的级别
//
func SetLevel(level LogLevel) {
	Level = level
}

//
// 设置日志打印的时区信息
//
func SetLocation(location *time.Location) {
	glocation = location
}

type myLogger struct {
	*log.Logger
	level    LogLevel
	codeinfo bool
	counter  *uint64
}

func newLogger(writer io.Writer, prefix string, level LogLevel, counter *uint64) *myLogger {
	return &myLogger{
		Logger:   log.New(writer, prefix, 0),
		level:    level,
		codeinfo: true,
		counter:  counter,
	}
}

func (l *myLogger) DisableCodeInfo() {
	l.codeinfo = false
}

func (l *myLogger) reallog(uniqId string, idx uint64, v ...interface{}) {
	tm := "[" + time.Now().In(glocation).Format("2006-01-02 15:04:05.999999") + "]"
	idx_s := fmt.Sprintf("[%d]", idx)
	if l.codeinfo && codeinfo {
		_, fn, line, _ := runtime.Caller(3)
		ss := strings.Split(fn, "/")
		sfn := ss[len(ss)-1]
		fline := fmt.Sprintf("%s:%d", sfn, line)
		v = append([]interface{}{tm, uniqId, idx_s, fline}, v...)
	} else {
		v = append([]interface{}{tm, uniqId, idx_s}, v...)
	}
	l.Logger.Println(v...)
}

func RequestId(ctx context.Context) string {
	i := ctx.Value(RequestIdKey)
	if i != nil {
		return i.(string)
	}
	return ""
}

func (l *myLogger) printlnX(ctx context.Context, v ...interface{}) {
	uqid := fmt.Sprintf("[reqid=%s]", RequestId(ctx))
	idx := atomic.AddUint64(l.counter, 1)
	l.reallog(uqid, idx, v...)
}

func Info(ctx context.Context, v ...interface{}) {
	glog.Info(ctx, v...)
}

func Debug(ctx context.Context, v ...interface{}) {
	glog.Debug(ctx, v...)
}

func Warning(ctx context.Context, v ...interface{}) {
	glog.Warning(ctx, v...)
}

func Error(ctx context.Context, v ...interface{}) {
	glog.Error(ctx, v...)
}

// 设置默认的Log
//
// 设置后，可以直接用 `rlib.Debug` 等打印日志
//
func SetDefaultLog(log *RLog) {
	glog = log
}

func DefaultLog() *RLog {
	return glog
}

var writeChan chan *lumberjack.Logger = make(chan *lumberjack.Logger, 10)

func care_rotate_log(writer *lumberjack.Logger) {
	writeChan <- writer
}

func run_rotate_log() {
	writes := []*lumberjack.Logger{}
	nextHandleTimer := next_run_interval()
	idleDelay := time.NewTimer(nextHandleTimer)
	for {
		select {
		case write := <-writeChan:
			writes = append(writes, write)
		case <-idleDelay.C:
			for _, write := range writes {
				write.Rotate()
			}
			time.Sleep(5 * time.Second)
			nextHandleTimer = next_run_interval()
			idleDelay.Reset(nextHandleTimer)
		}
	}
}

func next_run_interval() time.Duration {
	nowtime := time.Now().In(glocation)
	year, month, day := nowtime.Local().Date()
	nextHandle := time.Date(year, month, day, 23, 59, 59, 0, glocation)
	between := nextHandle.Sub(nowtime)
	if between.Seconds() < 0 {
		return time.Duration(float64(24*3600*time.Second)) + between
	}

	return time.Duration(between)
}
