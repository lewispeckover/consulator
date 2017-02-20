package main

import (
	"encoding/json"
	"flag"
	//"fmt"
	"github.com/fatih/color"
	"github.com/hashicorp/consul/api"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	path    = flag.String("path", "", "Path to file or directory containing data to load")
	debug   = flag.Bool("debug", false, "Show debugging information")
	trace   = flag.Bool("trace", false, "Show even more debugging information")
	dump    = flag.Bool("dump", false, "Dump loaded data as JSON, suitable for using in a `consul kv import`")
	enc     = json.NewEncoder(os.Stdout)
	absPath string
	data    api.KVPairs
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
	absPath, _ = filepath.Abs(*path)
	_, err := os.Stat(absPath)
	if err != nil {
		Error.Fatal(err)
	}
	err = filepath.Walk(absPath, parseConfig)
	if err != nil {
		Error.Fatal(err)
	}
	if *dump {
		exportData()
	}
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
		keyPrefix := strings.Split(
			// remove leading '/'
			strings.TrimPrefix(
				// remove the file extension
				strings.TrimSuffix(
					// remove the base path that was passed as -path
					strings.TrimPrefix(path, absPath),
					filepath.Ext(path)),
				string(os.PathSeparator)),
			string(os.PathSeparator))
		if keyPrefix[0] == "" {
			// remove the "" value if passed a file directly in -path
			keyPrefix = []string{}
		}
		Info.Printf("keyprefix is %v", keyPrefix)
		switch {
		case strings.HasSuffix(strings.ToLower(path), ".json"):
			Info.Printf("Parsing %s as json", path)
			err := parseJson(fp, keyPrefix)
			if err != nil {
				Warning.Printf("%v: %v\n", path, err)
			}
		case strings.HasSuffix(strings.ToLower(path), ".yml"):
			fallthrough
		case strings.HasSuffix(strings.ToLower(path), ".yaml"):
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
			err := parseJson(yamlR, keyPrefix)
			if err != nil {
				Warning.Printf("%v: %v\n", path, err)
			}
		case strings.HasSuffix(strings.ToLower(path), ".properties"):
			Info.Printf("Parsing %s as properties", path)
		case strings.HasSuffix(strings.ToLower(path), ".ini"):
			Info.Printf("Parsing %s as ini", path)
		default:
		}
	}
	return nil
}

func exportData() {
	exported := make([]*kvExportEntry, len(data))
	for i, kv := range data {
		exported[i] = toExportEntry(kv)
	}
	json, err := json.MarshalIndent(exported, "", "\t")
	if err != nil {
		Error.Printf("Error exporting data: %s", err)
		os.Exit(2)
	}
	os.Stdout.Write(json)
}
