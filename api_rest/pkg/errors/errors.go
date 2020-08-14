package errors

import "log"

func PanicError(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}
