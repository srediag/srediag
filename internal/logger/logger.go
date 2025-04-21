// internal/logger/logger.go
package logger

import "go.uber.org/zap"

var Sugar *zap.SugaredLogger

func Init(debug bool) error {
	var log *zap.Logger
	var err error
	if debug {
		log, err = zap.NewDevelopment()
	} else {
		log, err = zap.NewProduction()
	}
	if err != nil {
		return err
	}
	Sugar = log.Sugar()
	return nil
}
