package generr

import (
	"os"

	"github.com/sirupsen/logrus"
)

type LogParam struct {
	format       string
	output       string
	level        string
	indentation  bool
	showFileLine bool
}

func InitLog(format, output, level string, indent, file bool) (log *Log, err error) {
	serverLogParam := LogParam{
		format:       format,
		output:       output,
		level:        level,
		indentation:  indent,
		showFileLine: file,
	}
	return initLog(serverLogParam)
}

// New return a pointer to initiated logger
func New(format, output, level string, indent, file bool) (log *Log, err error) {
	serverLogParam := LogParam{
		format:       format,
		output:       output,
		level:        level,
		indentation:  indent,
		showFileLine: file,
	}

	return initLog(serverLogParam)
}

func initLog(p LogParam) (log *Log, err error) {
	log = &Log{}
	log.Logger = logrus.New()

	if p.showFileLine {
		hook := NewHook()
		hook.Field = "file"
		log.Logger.AddHook(hook)
	}

	setFormat(log.Logger, p)
	setOutput(log.Logger, p)
	err = setLevel(log.Logger, p)

	return log, err
}

func setFormat(log *logrus.Logger, p LogParam) {
	// TODO: should be completed
	switch p.format {
	case "json":
		log.SetFormatter(&logrus.JSONFormatter{
			PrettyPrint: p.indentation,
		})
	}
}

func setOutput(log *logrus.Logger, p LogParam) {
	switch p.output {
	case "stdout":
		log.SetOutput(os.Stdout)
	default:
		file, err := os.OpenFile(p.output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			log.Out = file
		} else {
			log.Fatalf("Failed to write logs to file %q, [***]", p.output)
		}
	}
}

func setLevel(log *logrus.Logger, p LogParam) (err error) {

	log.Level, err = logrus.ParseLevel(p.level)
	return
}
