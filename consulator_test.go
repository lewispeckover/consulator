package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestParseConfig(t *testing.T) {
	referenceData := map[string]string{
		"json/string":        "value",
		"json/int":           "1",
		"json/array":         "a,1,true",
		"json/bool":          "false",
		"json/nested/string": "value",
		"json/nested/int":    "1",
		"json/nested/array":  "a,1,true",
		"json/nested/bool":   "false",
		"yaml/string":        "value",
		"yaml/int":           "1",
		"yaml/array":         "a,1,true",
		"yaml/bool":          "false",
		"yaml/nested/string": "value",
		"yaml/nested/int":    "5294967294",
		"yaml/nested/array":  "a,1,true",
		"yaml/nested/bool":   "false",
	}
	nArgs = 1
	logInit()
	data = make(map[string][]byte)
	path = "./test"
	absPath, _ = filepath.Abs(path)
	*glue = ","

	_, err := os.Stat(absPath)
	if err != nil {
		t.Error(err)
	}
	err = filepath.Walk(absPath, parseConfig)
	if err != nil {
		t.Error(err)
	}
	for k, v := range referenceData {
		if val, ok := data[k]; ok {
			if string(val) == v {
				delete(referenceData, k)
				delete(data, k)
			}
		}
	}
	if len(referenceData) != 0 || len(data) != 0 {
		t.Error("Expected values do not match:")
		t.Error(" Unmatched data:\n")
		t.Error(fmt.Sprintf("  %#v\n", data))
		t.Error(" Unmatched reference data:\n")
		t.Error(fmt.Sprintf("  %#v\n", referenceData))
	}
}
