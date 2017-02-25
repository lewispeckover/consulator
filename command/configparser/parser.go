package configparser

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

var absPath string
var data map[string][]byte
var forceType = "auto"
var glue string

func Parse(path string, dataDest map[string][]byte, arrayGlue string) error {
	absPath, _ = filepath.Abs(path)
	data = dataDest
	glue = arrayGlue
	_, err := os.Stat(absPath)
	if err != nil {
		return err
	}
	err = filepath.Walk(absPath, walk)
	return err
}

func ParseAsJSON(path string, dataDest map[string][]byte, arrayGlue string) error {
	forceType = "json"
	return Parse(path, dataDest, arrayGlue)
}

func ParseAsYAML(path string, dataDest map[string][]byte, arrayGlue string) error {
	forceType = "yaml"
	return Parse(path, dataDest, arrayGlue)
}

func walk(path string, fstat os.FileInfo, err error) error {
	var fp *os.File
	if fstat.Mode().IsDir() {
		os.Stderr.WriteString("\nskipping dir ")
		os.Stderr.WriteString(path)
		return nil
	} else {
		fp, err = os.Open(path)
	}
	if err != nil {
		return err
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
		// remove blank token
		keyPrefix = []string{}
	}
	switch {
	case strings.HasSuffix(strings.ToLower(path), ".json") || forceType == "json":
		err := parseJson(fp, keyPrefix, glue)
		if err != nil {
			return err
		}
	case strings.HasSuffix(strings.ToLower(path), ".yml"):
		fallthrough
	case strings.HasSuffix(strings.ToLower(path), ".yaml") || forceType == "yaml":
		// yaml handling based on https://github.com/bronze1man/yaml2json
		yamlR, yamlW := io.Pipe()
		go func() {
			defer yamlW.Close()
			yamlToJson(fp, yamlW)
		}()
		err := parseJson(yamlR, keyPrefix, glue)
		if err != nil {
			return err
		}
	//case strings.HasSuffix(strings.ToLower(path), ".properties"):
	// TODO: .properties parsing
	//case strings.HasSuffix(strings.ToLower(path), ".ini"):
	// TODO: .ini parsing
	default:
	}
	return nil
}
