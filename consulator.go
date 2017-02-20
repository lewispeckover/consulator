package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/hashicorp/consul/api"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var usage = func() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] [PATH]\n\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "PATH should be the path to a file or directory that contains your data.")
	fmt.Fprintln(os.Stderr, "If no path is provided, stdin is used. In this case, -format must be specified.\n")
	fmt.Fprintln(os.Stderr, "Options:\n")
	flag.PrintDefaults()
}

var (
	debug   = flag.Bool("debug", false, "Show debugging information")
	dump    = flag.Bool("dump", false, "Dump loaded data as JSON, suitable for using in a 'consul kv import'")
	format  = flag.String("format", "", "Specify data format(json or yaml) when reading from stdin.")
	glue    = flag.String("glue", "\n", "Glue to use when joining array values")
	trace   = flag.Bool("trace", false, "Show even more debugging information")
	enc     = json.NewEncoder(os.Stdout)
	path    string
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
	flag.Usage = usage
	flag.Parse()
	logInit()
	switch flag.NArg() {
	case 0:
		// use stdin
		Trace.Println("No arguments, using stdin instead")
		fi, _ := os.Stdin.Stat()
		switch *format {
		case "json":
			fallthrough
		case "yaml":
			Trace.Printf("Stdin format is going to be: %v", *format)
			err := parseConfig(fmt.Sprintf(".%v", *format), fi, nil)
			if err != nil {
				Error.Printf("%v: %v\n", path, err)
			}
		default:
			Error.Fatal("When reading from stdin, the -format option must be provided and must one of: json, yaml")
		}
	case 1:
		path = flag.Arg(0)
		absPath, _ = filepath.Abs(path)
		_, err := os.Stat(absPath)
		if err != nil {
			Error.Fatal(err)
		}
		err = filepath.Walk(absPath, parseConfig)
		if err != nil {
			Error.Fatal(err)
		}
	default:
		Error.Printf("1 argument expected, but found %d\n\n", flag.NArg())
		usage()
		os.Exit(255)
	}
	if *dump {
		exportData()
	}
}

func parseConfig(path string, f os.FileInfo, err error) error {
	var fp *os.File
	Trace.Printf("Traversing %s", path)
	if f.Mode().IsRegular() && flag.NArg() == 1 {
		fp, err = os.Open(path)
	} else if flag.NArg() == 0 {
		fp = os.Stdin
	} else {
		fp, err = os.Open(os.DevNull)
		err = errors.New(fmt.Sprintf("Not a regular file: %s", path))
	}
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
	return nil
}

func exportData() {
	exported := make([]*kvExportEntry, len(data))
	for i, kv := range data {
		exported[i] = toExportEntry(kv)
	}
	json, err := json.MarshalIndent(exported, "", "\t")
	if err != nil {
		Error.Fatalf("Error exporting data: %s", err)
	}
	os.Stdout.Write(json)
	fmt.Println("")

}
