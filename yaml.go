package main

import (
	"encoding/json"
	"fmt"
	goyaml "gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"strconv"
)

// yaml handling based on https://github.com/bronze1man/yaml2json
func yamlToJson(in io.Reader, out io.Writer) error {
	input, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}
	var data interface{}
	err = goyaml.Unmarshal(input, &data)
	if err != nil {
		return err
	}
	input = nil
	err = yamlWalk(&data)
	if err != nil {
		return err
	}

	output, err := json.Marshal(data)
	if err != nil {
		return err
	}
	data = nil
	_, err = out.Write(output)
	return err
}

func yamlWalk(pIn *interface{}) (err error) {
	switch in := (*pIn).(type) {
	case map[interface{}]interface{}:
		m := make(map[string]interface{}, len(in))
		for k, v := range in {
			if err = yamlWalk(&v); err != nil {
				return err
			}
			var sk string
			switch k.(type) {
			case string:
				sk = k.(string)
			case int:
				sk = strconv.Itoa(k.(int))
			default:
				return fmt.Errorf("type mismatch: expect map key string or int; got: %T", k)
			}
			m[sk] = v
		}
		*pIn = m
	case []interface{}:
		for i := len(in) - 1; i >= 0; i-- {
			if err = yamlWalk(&in[i]); err != nil {
				return err
			}
		}
	}
	return nil
}
