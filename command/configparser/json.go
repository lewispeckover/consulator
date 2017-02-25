package configparser

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/antonholmquist/jason"
)

func parseJson(fp io.Reader, prefix []string, glue string) error {
	jsonObj, err := jason.NewObjectFromReader(fp)
	if err != nil {
		return err
	}
	j, err := jsonObj.GetObject()
	if err != nil {
		return err
	}
	return jsonWalk(prefix, j)
}

func jsonWalk(prefix []string, obj *jason.Object) error {
	for k, v := range obj.Map() {
		key := strings.Join(append(prefix, k), "/")
		switch v.Interface().(type) {
		case string:
			data[key] = []byte(fmt.Sprintf("%v", v.Interface()))
		case json.Number:
			data[key] = []byte(fmt.Sprintf("%v", v.Interface()))
		case []interface{}:
			// json array
			o, _ := v.Array()
			val, err := jsonArrayChoose(o)
			if err != nil {
				return err
			}
			data[key] = []byte(strings.Join(val, glue))
		case bool:
			data[key] = []byte(fmt.Sprintf("%v", v.Interface()))
		case nil:
			// json nulls
		case map[string]interface{}:
			// json object
			o, _ := v.Object()
			if err := jsonWalk(append(prefix, k), o); err != nil {
				return err
			}
		default:
		}
	}
	return nil
}

func jsonArrayChoose(arr []*jason.Value) (ret []string, err error) {
	for _, v := range arr {
		switch v.Interface().(type) {
		case string:
			ret = append(ret, fmt.Sprintf("%v", v.Interface()))
		case json.Number:
			ret = append(ret, fmt.Sprintf("%v", v.Interface()))
		case bool:
			ret = append(ret, fmt.Sprintf("%v", v.Interface()))
		default:
			return ret, fmt.Errorf(fmt.Sprintf("Invalid type %T in array. Only strings, numbers and boolean values are supported.\n", v.Interface()))
		}
	}
	return ret, nil
}
