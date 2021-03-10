package rlog

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

type LogFunc func(ctx context.Context, vals ...interface{})

type LogLevel int

const (
	LEVEL_DEBUG LogLevel = iota - 1
	LEVEL_INFO
	LEVEL_WARNING
	LEVEL_ERROR
)

const defName = "app"

var logs = sync.Map{}
var logLevel = LEVEL_INFO

func stdLogFunc(ctx context.Context, vals ...interface{}) {
	fmt.Fprintln(os.Stdout, vals...)
}

var (
	Debug   LogFunc = stdLogFunc
	Info    LogFunc = stdLogFunc
	Warning LogFunc = stdLogFunc
	Error   LogFunc = stdLogFunc
)

type Log struct {
	filename string
	logger   *log.Logger
	writer   *lumberjack.Logger
}

func GetLog(name string) *Log {
	raw, ok := logs.Load(name)
	if !ok {
		stdLogFunc(context.Background(), fmt.Sprintf("the log `%s` is not regist", name))
		return nil
	}
	return raw.(*Log)
}

func DefLog() *Log {
	return GetLog(defName)
}

func getLogFile(name string) string {
	return fmt.Sprintf("%s.%s.log", filepath.Base(os.Args[0]), name)
}

func New(name string) *Log {
	logfile := path.Join(config.Path, getLogFile(name))
	writer := &lumberjack.Logger{
		Filename:   logfile,
		MaxSize:    config.MaxMB,      // megabytes after which new file is created
		MaxBackups: config.MaxBackups, // number of backups
		MaxAge:     config.MaxDays,    //days
	}
	return &Log{
		filename: logfile,
		logger:   log.New(writer, "", 0),
		writer:   writer,
	}
}

func (g *Log) Debug(ctx context.Context, vals ...interface{}) {
	g.log(ctx, LEVEL_DEBUG, "Debug:", vals...)
}

func (g *Log) Info(ctx context.Context, vals ...interface{}) {
	g.log(ctx, LEVEL_INFO, "Info:", vals...)
}

func (g *Log) Warning(ctx context.Context, vals ...interface{}) {
	g.log(ctx, LEVEL_WARNING, "Warning:", vals...)
}

func (g *Log) Error(ctx context.Context, vals ...interface{}) {
	g.log(ctx, LEVEL_ERROR, "Error:", vals...)
}

func (g *Log) log(ctx context.Context, level LogLevel, prefix string, vals ...interface{}) {
	if logLevel > level {
		return
	}

	stime := "[" + time.Now().Format("2006-01-02 15:04:05.999999") + "]"

	_, fpath, fline, _ := runtime.Caller(3)
	fpaths := strings.Split(fpath, "/")
	fname := fpaths[len(fpaths)-1]
	codeinfo := fmt.Sprintf("%s:L%d", fname, fline)

	reqid := fmt.Sprintf("[req=%s]", CtxId(ctx))
	seq := fmt.Sprintf("[seq=%d]", logSerialNum(ctx))

	vals = append([]interface{}{prefix, stime, reqid, seq, codeinfo}, vals...)

	g.logger.Println(vals...)
}
