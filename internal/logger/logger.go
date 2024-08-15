package logger

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	errEmptyLogLevel       = errors.New("log level is empty")
	errInvalidLogFormatter = errors.New("invalid log formatter")
)

func Setup(level, formatter string) error {
	if level == "" {
		return errEmptyLogLevel
	}
	parsedLevel, err := logrus.ParseLevel(level)
	if err != nil {
		return fmt.Errorf("failed to parse log level: %w", err)
	}
	logOutput := strings.ToLower(formatter)
	switch logOutput {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	case "text":
		logrus.SetFormatter(&logrus.TextFormatter{})
	default:
		return fmt.Errorf("%w: failed to parse log formatter, not a valid logrus formatter '%s'", errInvalidLogFormatter, logOutput)
	}
	logrus.SetLevel(parsedLevel)
	return nil
}
