package helper

import "github.com/DataLabTechTV/labstore/backend/pkg/logger"

func Must[T any](val T, err error) T {
	if err != nil {
		logger.Log.Fatal(err)
	}
	return val
}

func CheckFatal(err error) {
	if err != nil {
		logger.Log.Fatal(err)
	}
}
