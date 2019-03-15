package service

import (
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// Log Variable
var log *logrus.Logger

// LogInit Function
func logInit() {
	// Initialize Log as New Logrus Logger
	log = logrus.New()

	// Set Log Format to JSON Format
	log.SetFormatter(&logrus.JSONFormatter{
		DisableTimestamp: false,
		TimestampFormat:  time.RFC3339Nano,
	})

	// Set Log Output to STDOUT
	log.SetOutput(os.Stdout)

	// Set Log Level
	switch strings.ToLower(os.Getenv("CONFIG_LOG_LEVEL")) {
	case "panic":
		log.SetLevel(logrus.PanicLevel)
	case "fatal":
		log.SetLevel(logrus.FatalLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "trace":
		log.SetLevel(logrus.TraceLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}
}

// Log Function
func Log(level string, label string, message string) {
	// Make Sure Log Is Not Empty Variable
	if log != nil {
		// Set Service Name Log Information
		service := strings.ToLower(os.Getenv("CONFIG_LOG_SERVICE"))

		// Print Log Based On Log Level Type
		switch strings.ToLower(level) {
		case "panic":
			log.WithFields(logrus.Fields{
				"service": service,
				"label":   label,
			}).Panic(message)
		case "fatal":
			log.WithFields(logrus.Fields{
				"service": service,
				"label":   label,
			}).Fatal(message)
		case "error":
			log.WithFields(logrus.Fields{
				"service": service,
				"label":   label,
			}).Error(message)
		case "warn":
			log.WithFields(logrus.Fields{
				"service": service,
				"label":   label,
			}).Warn(message)
		case "debug":
			log.WithFields(logrus.Fields{
				"service": service,
				"label":   label,
			}).Debug(message)
		case "tarce":
			log.WithFields(logrus.Fields{
				"service": service,
				"label":   label,
			}).Trace(message)
		default:
			log.WithFields(logrus.Fields{
				"service": service,
				"label":   label,
			}).Info(message)
		}
	}
}
