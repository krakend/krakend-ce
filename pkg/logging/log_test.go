package gologging

import (
	"bytes"
	"encoding/json"
	"io"
	"regexp"
	"strings"
	"testing"

	gologging "github.com/op/go-logging"
)

const (
	debugMsg    = "Debug msg"
	infoMsg     = "Info msg"
	warningMsg  = "Warning msg"
	errorMsg    = "Error msg"
	criticalMsg = "Critical msg"
)

func TestNewLogger(t *testing.T) {
	levels := []string{"DEBUG", "INFO", "WARNING", "ERROR", "CRITICAL"}
	regexps := []*regexp.Regexp{
		regexp.MustCompile(debugMsg),
		regexp.MustCompile(infoMsg),
		regexp.MustCompile(warningMsg),
		regexp.MustCompile(errorMsg),
		regexp.MustCompile(criticalMsg),
	}

	for i, level := range levels {
		output, err := logSomeStuff(level)
		if err != nil {
			t.Error(err)
			return
		}
		for j := i; j < len(regexps); j++ {
			if !regexps[j].MatchString(output) {
				t.Errorf("The output doesn't contain the expected msg for the level: %s. [%s]", level, output)
			}
		}
	}
}

func TestNewLogger_logstashFormat(t *testing.T) {
	buff := bytes.NewBuffer(make([]byte, 1024))
	SetFormatterSelector(func(w io.Writer) string {
		return ActivePattern
	})
	logger, err := NewLogger(newExtraConfig("DEBUG", "logstash", ""), TestFormatWriter{buff})
	if err != nil {
		t.Error(err)
		return
	}

	logger.Critical(criticalMsg)

	outputMsg := strings.ReplaceAll(buff.String(), "\x00", "")

	if !isJson(outputMsg) || !strings.HasPrefix(outputMsg, "{\"@timestamp\":") {
		t.Error("The output doesn't contain a logstash formatted log line")
	}
}

func TestNewLogger_customFormat(t *testing.T) {
	buff := bytes.NewBuffer(make([]byte, 1024))
	SetFormatterSelector(func(w io.Writer) string {
		return ActivePattern
	})
	logger, err := NewLogger(newExtraConfig("DEBUG", "custom", "----> %{message}"), TestFormatWriter{buff})
	if err != nil {
		t.Error(err)
		return
	}

	logger.Critical(criticalMsg)

	outputMsg := strings.ReplaceAll(buff.String(), "\x00", "")

	if outputMsg != "----> Critical msg\n" {
		t.Error("The output doesn't contain the custom format")
	}
}

func TestNewLogger_unknownLevel(t *testing.T) {
	_, err := NewLogger(newExtraConfig("UNKNOWN", "default", ""), bytes.NewBuffer(make([]byte, 1024)))
	if err == nil {
		t.Error("The factory didn't return the expected error")
		return
	}
	if err != gologging.ErrInvalidLogLevel {
		t.Errorf("The factory didn't return the expected error. Got: %s", err.Error())
	}
}

func newExtraConfig(level, format, customFormat string) map[string]interface{} {
	return map[string]interface{}{
		Namespace: map[string]interface{}{
			"level":         level,
			"prefix":        "pref",
			"syslog":        false,
			"stdout":        true,
			"format":        format,
			"custom_format": customFormat,
		},
	}
}

type TestFormatWriter struct {
	io.Writer
}

func logSomeStuff(level string) (string, error) {
	buff := bytes.NewBuffer(make([]byte, 1024))
	SetFormatterSelector(func(w io.Writer) string {
		switch w.(type) {
		case TestFormatWriter:
			return "customFormatter %{message}"
		default:
			return DefaultPattern
		}
	})
	logger, err := NewLogger(newExtraConfig(level, "default", ""), TestFormatWriter{buff})
	if err != nil {
		return "", err
	}

	logger.Debug(debugMsg)
	logger.Info(infoMsg)
	logger.Warning(warningMsg)
	logger.Error(errorMsg)
	logger.Critical(criticalMsg)

	return buff.String(), nil
}

func isJson(possibleJson string) bool {
	var js map[string]interface{}

	return json.Unmarshal([]byte(possibleJson), &js) == nil
}
