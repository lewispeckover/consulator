package main

import (
	"encoding/json"
	"fmt"
	"github.com/antonholmquist/jason"
	"github.com/hashicorp/consul/api"
	"io"
	"os"
	"strings"
)

func parseJson(fp io.Reader, prefix []string) error {
	jsonObj, err := jason.NewObjectFromReader(fp)
	if err != nil {
		Warning.Printf("%v: %v\n", path, err)
	}
	j, _ := jsonObj.GetObject()
	jsonWalk(prefix, j, err)
	return err
}

func jsonWalk(prefix []string, obj *jason.Object, err error) error {
	for k, v := range obj.Map() {
		key := strings.Join(append(prefix, k), "/")
		Trace.Printf("JSON iteration: %s", key)
		switch v.Interface().(type) {
		case string:
			Info.Printf("%v: \"%v\"\n", key, v.Interface())
			data = append(data, &api.KVPair{Key: key, Value: []byte(fmt.Sprintf("%v", v.Interface()))})
		case json.Number:
			Info.Printf("%v: \"%v\"\n", key, v.Interface())
			data = append(data, &api.KVPair{Key: key, Value: []byte(fmt.Sprintf("%v", v.Interface()))})
		case []interface{}:
			// json array
			o, _ := v.Array()
			Info.Printf("%v: \"%v\"\n", key, strings.Join(jsonArrayChoose(o), *glue))
			data = append(data, &api.KVPair{Key: key, Value: []byte(strings.Join(jsonArrayChoose(o), *glue))})
		case bool:
			Info.Printf("%v: \"%v\"\n", key, v.Interface())
			data = append(data, &api.KVPair{Key: key, Value: []byte(fmt.Sprintf("%v", v.Interface()))})
		case nil:
			// json nulls
		case map[string]interface{}:
			// json object
			o, _ := v.Object()
			jsonWalk(append(prefix, k), o, err)
		default:
			Warning.Printf("this is not a type we can handle: %T\n", v.Interface())
		}
	}
	return nil
}

func jsonArrayChoose(arr []*jason.Value) (ret []string) {
	for _, v := range arr {
		switch v.Interface().(type) {
		case string:
			ret = append(ret, fmt.Sprintf("%v", v.Interface()))
		case json.Number:
			ret = append(ret, fmt.Sprintf("%v", v.Interface()))
		case bool:
			ret = append(ret, fmt.Sprintf("%v", v.Interface()))
		default:
			Error.Printf("Invalid type %T in array. Only strings, numbers and boolean values are supported.\n", v.Interface())
			os.Exit(1)
		}
	}
	return ret
}
