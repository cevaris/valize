package valize

import (
	"log"
)

func FailOnErr(err error, msg string) {
	if err != nil {
		log.Fatal(err, msg)
	}
}

func LogOnErr(err error, msg string) {
	if err != nil {
		log.Fatal(err, msg)
	}
}
