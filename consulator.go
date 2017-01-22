package main

import (
  "os"
  "io"
  "flag"
  "fmt"
  "strings"
  "log"
  "path/filepath"
)

var path = flag.String("path", "", "Path to file or directory containing data to import")

func main() {
  flag.Parse()
  _, err := os.Stat(*path)
  if err != nil {
    log.Fatal(err)
  }
  err = filepath.Walk(*path, parseConfig)
  if err != nil {
    log.Fatal(err)
  }
}

func parseConfig(path string, f os.FileInfo, err error) error {
  if err != nil { log.Printf("%v: %v\n", path, err) }
  if f.Mode().IsRegular() {
    fp, err := os.Open(path)
    if err != nil { log.Printf("%v: %v\n", path, err) }
    switch {
      case strings.HasSuffix(path, ".json"):
        fmt.Println("json file found!")
        err := parseJson(fp)
        if err != nil { log.Printf("%v: %v\n", path, err) }
      case strings.HasSuffix(path, ".yml"):
        fallthrough
      case strings.HasSuffix(path, ".yaml"):
        fmt.Println("found a yaml file!")
        // yaml handling based on https://github.com/bronze1man/yaml2json
        yamlR, yamlW := io.Pipe()
        go func() {
          defer yamlW.Close()
          err := yamlToJson(fp, yamlW)
          if err != nil { log.Printf("%v: %v\n") }
        }()
        err := parseJson(yamlR)
        if err != nil { log.Printf("%v: %v\n", path, err) }
     case strings.HasSuffix(path, ".properties"):
        fmt.Println("found a properties file!")
      case strings.HasSuffix(path, ".ini"):
        fmt.Println("found an ini file!")
      default:
    }
  }
  return nil
}
