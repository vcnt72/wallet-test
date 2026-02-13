// Package logger
package logger

import "go.uber.org/zap"

var Log *zap.Logger

func Load() {
	z, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	Log = z
}
