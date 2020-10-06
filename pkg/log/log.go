package log

import (
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/dimaskiddo/go-whatsapp-rest/pkg/server"
)

// Log Variable
var logger *logrus.Logger

// Log Level Data Type
type logLevel string

// Log Level Data Type Constant
const (
	LogLevelPanic logLevel = "panic"
	LogLevelFatal logLevel = "fatal"
	LogLevelError logLevel = "error"
	LogLevelWarn  logLevel = "warn"
	LogLevelDebug logLevel = "debug"
	LogLevelTrace logLevel = "trace"
	LogLevelInfo  logLevel = "info"
)

// Initialize Function in Helper Logging
func init() {
	// Initialize Log as New Logrus Logger
	logger = logrus.New()

	// Set Log Format to JSON Format
	logger.SetFormatter(&logrus.JSONFormatter{
		DisableTimestamp: false,
		TimestampFormat:  time.RFC3339Nano,
	})

	// Set Log Output to STDOUT
	logger.SetOutput(os.Stdout)

	// Set Log Level
	switch strings.ToLower(server.Config.GetString("SERVER_LOG_LEVEL")) {
	case "panic":
		logger.SetLevel(logrus.PanicLevel)
	case "fatal":
		logger.SetLevel(logrus.FatalLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "trace":
		logger.SetLevel(logrus.TraceLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}
}

// Println Function
func Println(level logLevel, label string, message interface{}) {
	// Make Sure Log Is Not Empty Variable
	if logger != nil {
		// Set Service Name Log Information
		service := strings.ToLower(server.Config.GetString("SERVER_NAME"))

		// Print Log Based On Log Level Type
		switch level {
		case "panic":
			logger.WithFields(logrus.Fields{
				"service": service,
				"label":   label,
			}).Panicln(message)
		case "fatal":
			logger.WithFields(logrus.Fields{
				"service": service,
				"label":   label,
			}).Fatalln(message)
		case "error":
			logger.WithFields(logrus.Fields{
				"service": service,
				"label":   label,
			}).Errorln(message)
		case "warn":
			logger.WithFields(logrus.Fields{
				"service": service,
				"label":   label,
			}).Warnln(message)
		case "debug":
			logger.WithFields(logrus.Fields{
				"service": service,
				"label":   label,
			}).Debug(message)
		case "trace":
			logger.WithFields(logrus.Fields{
				"service": service,
				"label":   label,
			}).Traceln(message)
		default:
			logger.WithFields(logrus.Fields{
				"service": service,
				"label":   label,
			}).Infoln(message)
		}
	}
}
