package main

import (
	"encoding/base64"
	"fmt"
	"github.com/hashicorp/consul/api"
	//"strings"
)

var config = api.DefaultConfig()
var client, _ = api.NewClient(config)

func dumpConsul() {
	fmt.Println("dumping consul kv")
	kv := client.KV()
	keys, _, err := kv.Keys("/", "/", &api.QueryOptions{})
	if err != nil {
		panic(err)
	}
	for _, key := range keys {
		fmt.Printf("%s\n", key)
	}
}

type kvExportEntry struct {
	Key   string `json:"key"`
	Flags uint64 `json:"flags"`
	Value string `json:"value"`
}

func toExportEntry(key string, val []byte) *kvExportEntry {
	return &kvExportEntry{
		Key:   key,
		Flags: 0,
		Value: base64.StdEncoding.EncodeToString(val),
	}
}
