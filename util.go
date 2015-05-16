package main

import (
	"log"
	"runtime"
)

func callerName() string {
	pc, _, _, _ := runtime.Caller(1)
	return "[" + runtime.FuncForPC(pc).Name() + "]"
}

func fatal(err error) {
	if err != nil {
		log.Fatalln(callerName(), err)
	}
}

func check(err error) {
	if err != nil {
		log.Println(callerName(), err)
	}
}
