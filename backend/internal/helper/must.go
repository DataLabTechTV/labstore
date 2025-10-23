package helper

import log "github.com/sirupsen/logrus"

func Must[T any](val T, err error) T {
	if err != nil {
		log.Fatal(err)
	}
	return val
}

func MustGet[T any](val T, ok bool, msg string) T {
	if !ok {
		log.Fatal(msg)
	}
	return val
}

func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
