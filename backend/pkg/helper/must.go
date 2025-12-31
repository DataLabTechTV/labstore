package helper

import (
	"fmt"
	"os"
)

func Must[T any](val T, err error) T {
	CheckFatal(err)
	return val
}

func CheckFatal(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
