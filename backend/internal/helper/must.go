package helper

func Must[T any](val T, err error) T {
	CheckFatal(err)
	return val
}

func CheckFatal(err error) {
	if err != nil {
		panic(err)
	}
}
