package logging

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// Message holds the existing log message fields
type Message struct {
	Severity string
	Tag      string
	Body     interface{}
	Location string
	Debug    bool
}

type logFormatter struct {
	logrus.TextFormatter
}

// Formats log output with timestamp, colors and proper severity
func (f *logFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = 35 // purple
	case logrus.WarnLevel:
		levelColor = 33 // yellow
	case logrus.ErrorLevel:
		levelColor = 31 // red
	case logrus.FatalLevel, logrus.PanicLevel:
		levelColor = 41 // white letters + red bg
	default:
		levelColor = 36 // blue
	}
	return []byte(
		fmt.Sprintf(
			"%s [\x1b[%dm%s\x1b[0m] %s\n",
			entry.Time.Format(f.TimestampFormat),
			levelColor,
			strings.ToUpper(entry.Level.String()),
			entry.Message,
		)), nil
}

// Info will log informational messages to stdout and logs
func Info(tag string, message interface{}) {
	writeMessage(Message{Severity: "INFO", Tag: tag, Body: message})
}

// Error will log Error messages to stdout and logs and provide information on runtime caller,
// preferably through 'logging.CallerInfo()' function
func Error(tag string, err interface{}, callerInfo string) {
	writeMessage(Message{Severity: "ERROR", Tag: tag, Body: err, Location: callerInfo})
}

func writeMessage(m Message) {
	file := "restaurantsearch_" + time.Now().Format("2006-01-02") + ".log"
	f, er1 := os.OpenFile(filepath.Clean(file), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if er1 != nil {
		log.Fatal("Failed to write to log file. Exiting!")
	}

	logger := &logrus.Logger{
		Out:   io.MultiWriter(os.Stderr, f),
		Level: logrus.InfoLevel,
		Formatter: &logFormatter{
			logrus.TextFormatter{
				FullTimestamp:          true,
				TimestampFormat:        "2006-01-02 15:04:05",
				ForceColors:            true,
				DisableLevelTruncation: true,
			},
		},
	}

	ctx := context.Background()
	entry := fmt.Sprintf("[%v] %v", m.Tag, m.Body)

	taggedFields := logrus.Fields{"severity": m.Severity, "tag": m.Tag, "body": m.Body, "caller": m.Location}

	switch m.Severity {
	case "INFO":
		logger.WithFields(taggedFields).WithContext(ctx).Info(entry)
	case "ERROR":
		entry := fmt.Sprintf("%v [%v]", entry, m.Location)
		logger.WithFields(taggedFields).WithContext(ctx).Error(entry)
	default:
		logger.WithFields(taggedFields).WithContext(ctx).Info(entry)
	}

	er2 := f.Close()
	if er2 != nil {
		log.Fatal(er2)
	}
}

// CallerInfo returns function name, file and line number which invoked it
func CallerInfo(depthList ...int) string {
	var depth int
	if depthList == nil {
		depth = 1
	} else {
		depth = depthList[0]
	}
	function, file, line, _ := runtime.Caller(depth)
	return fmt.Sprintf("%s.%s:%d", filepath.Base(file), filepath.Base(runtime.FuncForPC(function).Name()), line)
}
