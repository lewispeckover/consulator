package main

import (
  "io"
  "fmt"
  "strings"
  "encoding/json"
  "github.com/antonholmquist/jason"
)

func parseJson(fp io.Reader) (error){
  jsonObj, err := jason.NewObjectFromReader(fp)
  if err != nil { fmt.Printf("%v: %v\n", path, err) }
  j, _ := jsonObj.GetObject()
  jsonWalk("", j, err)
  return err
}

func jsonWalk(prefix string, obj *jason.Object, err error) (error) {
  for k, v := range obj.Map() {
    switch v.Interface().(type) {
      case string:
        fmt.Printf("%v: %v\n", prefix + "/" + k, v.Interface())
        enc.Encode(v.Interface())
      case json.Number:
        fmt.Printf("%v: %v\n", prefix + "/" + k, v.Interface())
        enc.Encode(v.Interface())
      case []interface {}:
        // json array
        o, _ := v.Array()
        fmt.Printf("%v: %v\n", prefix + "/" + k, strings.Join(jsonArrayChoose(o), ", "))
        enc.Encode(v.Interface())
      case map[string]interface {}:
        // json object
        o, _ := v.Object()
        jsonWalk(prefix + "/" + k, o, err)
      case bool:
        fmt.Printf("%v: %v\n", prefix + "/" + k, v.Interface())
        enc.Encode(v.Interface())
      case nil:
        // json nulls 
      default:
        fmt.Printf("this is not a type we can handle: %T\n", v.Interface())
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
        fmt.Printf("Ignoring type %T in array. Only strings, numbers and boolean values are supported.\n", v.Interface())
    }
  }
  return ret
}
