package utils

import "log"

func Checkerr(err error) {
	if err != nil {
		log.Fatalf(err.Error())
	}
}
