package core

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"path"
)

type MyLog struct {
}

// Colors for different levels
const (
	red    = 31
	yellow = 33
	blue   = 36
	gray   = 37
)

func (MyLog) Format(entry *logrus.Entry) ([]byte, error) {
	// Color
	var color int
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		color = gray
	case logrus.WarnLevel:
		color = yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		color = red
	default:
		color = blue
	}

	// Buffer is required for formatting log messages before outputting them.
	var buf *bytes.Buffer
	if entry.Buffer != nil {
		buf = entry.Buffer
	} else {
		buf = &bytes.Buffer{}
	}

	// Time format
	timeFormat := entry.Time.Format("2006-01-02T15:04:05Z0700")

	if entry.HasCaller() {
		// Custom file path and line
		funcVal := entry.Caller.Function
		fileVal := fmt.Sprintf("%s:%d", path.Base(entry.Caller.File), entry.Caller.Line)
		// Custom format
		_, err := fmt.Fprintf(buf, "[%s] \x1b[%dm[%s]\x1b[0m %s %s %s\n", timeFormat, color, entry.Level, fileVal, funcVal, entry.Message)
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func InitLogger() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(MyLog{})
	//logrus.SetFormatter(&logrus.JSONFormatter{})
}
