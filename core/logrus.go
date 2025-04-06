package core

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"sync"
)

// Colors for different levels
const (
	red    = 31
	yellow = 33
	blue   = 36
	gray   = 37
)

type MyLog struct {
}

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

type MyHook struct {
	file     *os.File   // Log file
	errFile  *os.File   // Error log file
	fileDate string     // Date of log file
	logPath  string     // Path of log file
	mu       sync.Mutex // Mutex lock
}

func (hook *MyHook) Fire(entry *logrus.Entry) error {
	hook.mu.Lock()
	defer hook.mu.Unlock()

	date := entry.Time.Format("2006-01-02")
	if hook.fileDate != date {
		// Rotate if day is passed
		if err := hook.rotate(date); err != nil {
			return err
		}
	}

	// Dump logs to file
	entryStr, err := entry.String()
	if err != nil {
		return fmt.Errorf("failed to get log entry: %v", err)
	}
	if _, err := hook.file.Write([]byte(entryStr)); err != nil {
		return fmt.Errorf("failed to write to log file: %v", err)
	}

	// Dump error logs to file
	if entry.Level <= logrus.ErrorLevel {
		if _, err := hook.errFile.Write([]byte(entryStr)); err != nil {
			return fmt.Errorf("failed to write to error log file: %v", err)
		}
	}

	return nil
}

func (hook *MyHook) rotate(date string) error {
	if hook.file != nil {
		// Close the old one
		if err := hook.file.Close(); err != nil {
			return fmt.Errorf("failed to close the old log file when rotation: %v", err)
		}
	}
	if hook.errFile != nil {
		// Close the old one
		if err := hook.errFile.Close(); err != nil {
			return fmt.Errorf("failed to close the old error log file when rotation: %v", err)
		}
	}

	// Log file directory
	dir := fmt.Sprintf("%s/%s", hook.logPath, date)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create log directory: %v", err)
	}

	infoLog := fmt.Sprintf("%s/info.log", dir)
	errLog := fmt.Sprintf("%s/err.log", dir)

	// Create new log files
	var err error
	hook.file, err = os.OpenFile(infoLog, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}
	hook.errFile, err = os.OpenFile(errLog, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("failed to open error log file: %v", err)
	}

	// Update file date
	hook.fileDate = date
	return nil
}

func (hook *MyHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func InitLogger() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(MyLog{})
	//logrus.SetFormatter(&logrus.JSONFormatter{})

	logrus.AddHook(&MyHook{
		logPath: "logs",
	})
}
