package main

import (
	"github.com/fatih/color"
	"io/ioutil"
	"log"
	"os"
)

func logInit() {
	if *trace {
		Trace = log.New(os.Stderr, "TRACE: ", 0)
		Info = log.New(os.Stderr, color.BlueString("INFO: "), 0)
	} else if *debug {
		Trace = log.New(ioutil.Discard, "TRACE: ", 0)
		Info = log.New(os.Stderr, color.BlueString("INFO: "), 0)

	} else {
		Trace = log.New(ioutil.Discard, "TRACE: ", 0)
		Info = log.New(ioutil.Discard, color.BlueString("INFO: "), 0)
	}
	Warning = log.New(os.Stderr, color.RedString("WARNING: "), 0)
	Error = log.New(os.Stderr, color.RedString("ERROR: "), 0)
}
