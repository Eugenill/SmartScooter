package errors

func PanicError(err error) {
	if err != nil {
		panic(err)
	}
}
