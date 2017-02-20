package main

import (
	"encoding/json"
	"flag"
	//"fmt"
	"github.com/fatih/color"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	path    = flag.String("path", "", "Path to file or directory containing data to import")
	debug   = flag.Bool("debug", false, "Show debugging information")
	trace   = flag.Bool("trace", false, "Show even more debugging information")
	enc     = json.NewEncoder(os.Stdout)
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
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

func main() {
	flag.Parse()
	logInit()
	_, err := os.Stat(*path)
	if err != nil {
		Error.Fatal(err)
	}
	err = filepath.Walk(*path, parseConfig)
	if err != nil {
		Error.Fatal(err)
	}
	Dump()
}

func parseConfig(path string, f os.FileInfo, err error) error {
	Trace.Printf("Traversing %s", path)
	if err != nil {
		Warning.Printf("%v: %v\n", path, err)
	}
	if f.Mode().IsRegular() {
		fp, err := os.Open(path)
		if err != nil {
			Warning.Printf("%v: %v\n", path, err)
		}
		switch {
		case strings.HasSuffix(path, ".json"):
			Info.Printf("Parsing %s as json", path)
			err := parseJson(fp)
			if err != nil {
				Warning.Printf("%v: %v\n", path, err)
			}
		case strings.HasSuffix(path, ".yml"):
			fallthrough
		case strings.HasSuffix(path, ".yaml"):
			Info.Printf("Parsing %s as yaml", path)
			// yaml handling based on https://github.com/bronze1man/yaml2json
			yamlR, yamlW := io.Pipe()
			go func() {
				defer yamlW.Close()
				err := yamlToJson(fp, yamlW)
				if err != nil {
					Warning.Printf("%v: %v\n", path, err)
				}
			}()
			err := parseJson(yamlR)
			if err != nil {
				Warning.Printf("%v: %v\n", path, err)
			}
		case strings.HasSuffix(path, ".properties"):
			Info.Printf("Parsing %s as properties", path)
		case strings.HasSuffix(path, ".ini"):
			Info.Printf("Parsing %s as ini", path)
		default:
		}
	}
	return nil
}
