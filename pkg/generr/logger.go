package generr

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
)

type Log struct {
	Logger *logrus.Logger
}

// CheckError print all errors which happened inside the services, mainly they just have an error and a message
func (log *Log) CheckError(err error, code, message string, data ...interface{}) {
	message = fmt.Sprintf("CODE: %v , MESSAGE: %v", code, message)
	if err != nil {
		log.LogError(err, message, data...)
	}
}

// LogError record an error with message and variadic parameters
func (log *Log) LogError(err error, message string, data ...interface{}) {
	if err == nil && message != "" {
		err = fmt.Errorf(message)
	}
	if err != nil {
		if data == nil {
			log.Logger.WithFields(logrus.Fields{
				"err": err.Error(),
			}).Error(message)
		} else {
			log.Logger.WithFields(logrus.Fields{
				"err":  err.Error(),
				"data": fmt.Sprintf("%+v", data),
			}).Error(message)

		}
	}
}

// Debug print struct with details with logrus ability
func (log *Log) Debug(objs ...interface{}) {
	for _, v := range objs {
		parts := make(map[string]interface{}, 2)
		parts["type"] = fmt.Sprintf("%T", v)
		parts["value"] = v
		dataInJSON, _ := json.Marshal(parts)

		log.Logger.Debug(string(dataInJSON))
	}
}

// CheckInfo record the info
func (log *Log) CheckInfo(err error, message string, data ...interface{}) {
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Info(message)
		if data != nil {
			log.Debug(data...)
		}
	}
}

// Info is information
func (log *Log) Info(data ...interface{}) {
	log.Logger.Info(data...)
}

// Error is error
func (log *Log) Error(data ...interface{}) {
	log.Logger.Error(data...)
}

// Fatal stop application
func (log *Log) Fatal(data ...interface{}) {
	log.Logger.Fatal(data...)
}
