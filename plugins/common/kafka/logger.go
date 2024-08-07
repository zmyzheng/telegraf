package kafka

import (
	"github.com/IBM/sarama"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/logger"
)

type Logger struct {
}

// DebugLogger logs messages from sarama at the debug level.
type DebugLogger struct {
	Log telegraf.Logger
}

func (l *DebugLogger) Print(v ...interface{}) {
	l.Log.Debug(v...)
}

func (l *DebugLogger) Printf(format string, v ...interface{}) {
	l.Log.Debugf(format, v...)
}

func (l *DebugLogger) Println(v ...interface{}) {
	l.Print(v...)
}

// SetLogger configures a debug logger for kafka (sarama)
func (k *Logger) SetLogger() {
	sarama.Logger = &DebugLogger{Log: logger.New("sarama", "", "")}
}
