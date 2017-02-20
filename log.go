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
		Debug = log.New(os.Stderr, "DEBUG: ", 0)
	} else if *debug {
		Trace = log.New(ioutil.Discard, "TRACE: ", 0)
		Debug = log.New(os.Stderr, "DEBUG: ", 0)

	} else {
		Trace = log.New(ioutil.Discard, "TRACE: ", 0)
		Debug = log.New(ioutil.Discard, "DEBUG: ", 0)
	}
	if *quiet {
		Info = log.New(ioutil.Discard, color.BlueString("INFO: "), 0)
		Warning = log.New(ioutil.Discard, color.MagentaString("WARNING: "), 0)
	} else {
		Info = log.New(os.Stderr, color.BlueString("INFO: "), 0)
		Warning = log.New(os.Stderr, color.MagentaString("WARNING: "), 0)
	}
	Error = log.New(os.Stderr, color.RedString("ERROR: "), 0)
}
