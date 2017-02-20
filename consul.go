package main

import (
	//"encoding/base64"
	"fmt"
	"github.com/hashicorp/consul/api"
	//"strings"
)

var config = api.DefaultConfig()
var client, _ = api.NewClient(config)

func Dump() {
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
