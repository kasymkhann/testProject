package logging

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

// writerHook хук для направления записей журнала в несколько writer'ов в зависимости от уровней логирования.
type writerHook struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

// Fire отправляет записи журнала в указанные writer'ы в зависимости от уровней логирования.
func (hook *writerHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	for _, w := range hook.Writer {
		w.Write([]byte(line))
	}
	return err
}

// Levels возвращает уровни логирования, которые активируют хук.
func (hook *writerHook) Levels() []logrus.Level {
	return hook.LogLevels
}

var e *logrus.Entry

// Структурированный логгер, расширяющий тип Entry из Logrus.
type Logger struct {
	*logrus.Entry
}

// GetLogger создает новый экземпляр логгера.
func GetLogger() *Logger {
	return &Logger{e}
}

// GetLoggerWithField создает новый экземпляр логгера с дополнительным полем.
func (l *Logger) GetLoggerWithField(k string, v interface{}) *Logger {
	return &Logger{l.WithField(k, v)}
}

// init Инициализирует логгер с указанными конфигурациями.
func init() {
	l := logrus.New()
	l.SetReportCaller(true)
	l.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			filename := path.Base(frame.File)
			return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s:%d", filename, frame.Line)
		},
		DisableColors: false,
		FullTimestamp: true,
	}

	err := os.MkdirAll("logs", 0755) // 0644
	if err != nil {
		panic(err)
	}

	allFile, err := os.OpenFile("logs/all.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		logrus.Fatal(err)
	}

	l.SetOutput(io.Discard)

	l.AddHook(&writerHook{
		Writer:    []io.Writer{allFile, os.Stdout},
		LogLevels: logrus.AllLevels,
	})

	l.SetLevel(logrus.TraceLevel)

	e = logrus.NewEntry(l)
}
