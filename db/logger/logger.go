package logger

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/eachinchung/log"
	gormlogger "gorm.io/gorm/logger"
)

const (
	Reset       = "\033[0m"
	Red         = "\033[31m"
	Green       = "\033[32m"
	Yellow      = "\033[33m"
	Magenta     = "\033[35m"
	BlueBold    = "\033[34;1m"
	MagentaBold = "\033[35;1m"
	RedBold     = "\033[31;1m"
)

// Define gorm log levels.
const (
	Silent gormlogger.LogLevel = iota + 1
	Error
	Warn
	Info
)

// Writer log writer interface.
type Writer interface {
	Printf(string, ...any)
}

// Config defines a gorm logger configuration.
type Config struct {
	SlowThreshold time.Duration
	Colorful      bool
	LogLevel      gormlogger.LogLevel
}

type logger struct {
	Writer
	Config
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

// New create a gorm logger instance.
func New(level int, enableColor bool) gormlogger.Interface {
	var (
		infoStr      = "%s[info] "
		warnStr      = "%s[warn] "
		errStr       = "%s[error] "
		traceStr     = "[%s][%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s[%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s[%.3fms] [rows:%v] %s"
	)

	config := Config{
		SlowThreshold: 200 * time.Millisecond,
		Colorful:      enableColor,
		LogLevel:      gormlogger.LogLevel(level),
	}

	if config.Colorful {
		infoStr = Green + "%s " + Reset + Green + "[info] " + Reset
		warnStr = BlueBold + "%s " + Reset + Magenta + "[warn] " + Reset
		errStr = Magenta + "%s " + Reset + Red + "[error] " + Reset
		traceStr = Green + "%s " + Reset + Yellow + "[%.3fms] " + BlueBold + "[rows:%v]" + Reset + " %s"
		traceWarnStr = Green + "%s " + Yellow + "%s " + Reset + RedBold + "[%.3fms] " + Yellow + "[rows:%v]" + Magenta + " %s" + Reset
		traceErrStr = RedBold + "%s " + MagentaBold + "%s " + Reset + Yellow + "[%.3fms] " + BlueBold + "[rows:%v]" + Reset + " %s"
	}

	return &logger{
		Writer:       log.StdInfoLogger(),
		Config:       config,
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
}

// LogMode log mode.
func (l *logger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.LogLevel = level

	return &newLogger
}

// Info print info.
func (l logger) Info(ctx context.Context, msg string, data ...any) {
	if l.LogLevel >= Info {
		l.Printf(l.infoStr+msg, append([]any{fileWithLineNum()}, data...)...)
	}
}

// Warn print warn messages.
func (l logger) Warn(ctx context.Context, msg string, data ...any) {
	if l.LogLevel >= Warn {
		l.Printf(l.warnStr+msg, append([]any{fileWithLineNum()}, data...)...)
	}
}

// Error print error messages.
func (l logger) Error(ctx context.Context, msg string, data ...any) {
	if l.LogLevel >= Error {
		l.Printf(l.errStr+msg, append([]any{fileWithLineNum()}, data...)...)
	}
}

// Trace print sql message.
func (l logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= 0 {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= Error:
		sql, rows := fc()
		if rows == -1 {
			l.Printf(l.traceErrStr, fileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Printf(l.traceErrStr, fileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			l.Printf(l.traceWarnStr, fileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Printf(l.traceWarnStr, fileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case l.LogLevel >= Info:
		sql, rows := fc()
		if rows == -1 {
			l.Printf(l.traceStr, fileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Printf(l.traceStr, fileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}

func fileWithLineNum() string {
	for i := 4; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)

		if ok && !strings.HasSuffix(file, "_test.go") {
			dir, f := filepath.Split(file)

			return filepath.Join(filepath.Base(dir), f) + ":" + strconv.FormatInt(int64(line), 10)
		}
	}

	return ""
}
